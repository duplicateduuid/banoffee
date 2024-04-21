import heapq
import os
from urllib.parse import parse_qs, urlparse

import psycopg2
from psycopg2 import pool

from surprise import Dataset, Reader, KNNBaseline # type: ignore

from http.server import BaseHTTPRequestHandler, HTTPServer
import json

def load_config():
    import configparser
    config = configparser.ConfigParser()
    config.read("database.ini")
    return {
        "host": config.get("postgresql", "host"),
        "database": config.get("postgresql", "database"),
        "user": config.get("postgresql", "user"),
        "password": config.get("postgresql", "password")
    }

def create_database_pool():
    config = load_config()

    pool = psycopg2.pool.SimpleConnectionPool(
        minconn=1,
        maxconn=20,
        **config
    )

    return pool

def generate_datasets(pool):
    connection = pool.getconn()

    if not os.path.exists('./datasets'):
        os.mkdir('./datasets')

    with open("./datasets/resources.csv", "w+") as csv_file:
        cursor = connection.cursor()
        cursor.copy_expert("COPY (SELECT id, url, name, author FROM resource) TO STDOUT WITH CSV HEADER", csv_file)

    with open("./datasets/ratings.csv", "w+") as csv_file:
        cursor = connection.cursor()
        cursor.copy_expert("""
        COPY (
            SELECT 
                user_id AS user, 
                resource_id AS item, 
                CASE
                    WHEN review_rating IS NULL THEN 3
                    WHEN review_rating = 'one' THEN 1
                    WHEN review_rating = 'two' THEN 2
                    WHEN review_rating = 'three' THEN 3
                    WHEN review_rating = 'four' THEN 4
                    WHEN review_rating = 'five' THEN 5
                END AS rating
            FROM 
                user_resource
        ) TO STDOUT WITH CSV HEADER
        """, csv_file)

    pool.putconn(connection)

def get_resource(connection, url):
    cur = connection.cursor()
    sql_query = """
        SELECT 
            r.id, r.url, r.name, r.image_url, r.author, r.description
        FROM 
            "resource" r 
        WHERE 
            r.url=%s
    """
    cur.execute(sql_query, (url,))
    result = cur.fetchone()

    result_dict = {
        'id': result[0],
        'url': result[1],
        'name': result[2],
        'image_url': result[3],
        'author': result[4],
        'description': result[5]
    }

    return result_dict

class CustomKNN(KNNBaseline):
    def custom_get_neighbors(self, item_id, user_id, exclude_list, k):
        user_ratings = {iid for (iid, _) in self.trainset.ur[self.trainset.to_inner_uid(user_id)]}
        others = [(iid, self.sim[item_id, iid]) for iid in self.trainset.all_items() if iid not in user_ratings and iid not in exclude_list]
        k_nearest_neighbors = heapq.nlargest(k, others, key=lambda tple: tple[1])

        return k_nearest_neighbors
 
def train_knn(algo, pool):
    generate_datasets(pool)

    id_to_url, url_to_id = read_resources_urls()

    file_path = os.path.expanduser("./datasets/ratings.csv")
    reader = Reader(line_format="user item rating", sep=",", skip_lines=1)

    data = Dataset.load_from_file(file_path, reader=reader)

    trainset = data.build_full_trainset()
    algo.fit(trainset)

    return id_to_url, url_to_id

def read_resources_urls():
    file_name = "./datasets/resources.csv"
    id_to_url = {}
    url_to_id = {}
    with open(file_name, encoding="ISO-8859-1") as f:
        first_line = True

        for line in f:
            if first_line:
                first_line = False
                continue
        
            line = line.split(',')
            id_to_url[line[0]] = line[1]
            url_to_id[line[1]] = line[0]
    
    return id_to_url, url_to_id

def get_resource_nearest_neighbors(algo, id_to_url, url_to_id, resource_url, user_id, exclude_url_list, k=10):
    raw_id = url_to_id[resource_url]
    inner_id = algo.trainset.to_inner_iid(raw_id)

    exclude_list = [algo.trainset.to_inner_iid(url_to_id[url]) for (url, _) in exclude_url_list]
    neighbors = algo.custom_get_neighbors(inner_id, user_id, exclude_list, k)
    neighbors = (
        (algo.trainset.to_raw_iid(inner_id), score) for (inner_id, score) in neighbors
    )
    neighbors = ((id_to_url[rid], score) for (rid,score) in neighbors)

    return list(neighbors)

def calculate_items_per_query(query_results_length, n):
    max_items_per_query = min(n, 5)
    items_per_query = min(max_items_per_query, max(1, n // query_results_length))
    
    return items_per_query

class HTTPRequestHandler(BaseHTTPRequestHandler):
    def __init__(self, pool, algo, id_to_url, url_to_id, *args, **kwargs):
        self.pool = pool
        self.algo = algo
        self.id_to_url = id_to_url
        self.url_to_id = url_to_id
        super().__init__(*args, **kwargs)

    def do_POST(self):
        try:
            train_knn(self.algo, self.pool)

            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            response_json = json.dumps({ 'message': 'algo updated with success!' })
            self.wfile.write(response_json.encode('utf-8'))
        except:
            self.send_response(500)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            response_json = json.dumps({ 'message': 'failed to update algo!' })
            self.wfile.write(response_json.encode('utf-8'))

    # TODO: do some sort of pagination with offset here
    def do_GET(self):
        parsed_url = urlparse(self.path)
        query_params = parse_qs(parsed_url.query)

        user_id = parsed_url.path.strip('/')
        limit = int(query_params.get('limit', [''])[0] or 10)

        connection = self.pool.getconn()
        
        try:
            cur = connection.cursor()
            sql_query = sql_query = """
                SELECT 
                    r.url
                FROM 
                    user_resource ur
                LEFT JOIN 
                    resource r ON r.id = ur.resource_id
                WHERE 
                    ur.user_id = %s
                ORDER BY ur.updated_at DESC
                LIMIT 5;
            """
            cur.execute(sql_query, (user_id,))
            results = cur.fetchall()

            if len(results) <= 0:
                self.send_response(200)
                self.send_header('Content-Type', 'application/json')
                self.end_headers()
                response_json = json.dumps({ 'recommendations': [] })
                self.wfile.write(response_json.encode('utf-8'))
                return

            recommendations = []

            for row in results:
                nearest_neighbors = get_resource_nearest_neighbors(algo, self.id_to_url, self.url_to_id, row[0], user_id, recommendations, calculate_items_per_query(len(results), limit))
                recommendations.extend(nearest_neighbors)
            
            # Sort recommendations by score
            recommendations = heapq.nlargest(len(recommendations), recommendations, key=lambda tple: tple[1])
            recommendations = [get_resource(connection, url) for (url, _) in recommendations]

            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            response_json = json.dumps({ 'recommendations': recommendations })
            self.wfile.write(response_json.encode('utf-8'))
        except Exception as error:
            print(f"An exception occurred: {repr(error)} while trying to get recommendations for the user: {user_id}")

            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            response_json = json.dumps({ 'recommendations': [] })
            self.wfile.write(response_json.encode('utf-8'))
        finally:
            self.pool.putconn(connection)


def run_server(pool, algo, id_to_url, url_to_id, port=8000):
    server_address = ('', port)
    httpd = HTTPServer(server_address, lambda *args, **kwargs: HTTPRequestHandler(pool, algo, id_to_url, url_to_id, *args, **kwargs))
    print(f'Starting server on port {port}...')
    httpd.serve_forever()

pool = create_database_pool()

sim_options = {"name": "cosine", "user_based": False}
algo = CustomKNN(sim_options=sim_options)
id_to_url, url_to_id = train_knn(algo, pool)

run_server(pool, algo, id_to_url, url_to_id)