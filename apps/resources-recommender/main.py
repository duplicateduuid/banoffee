import os

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

def train_knn(algo, pool):
    generate_datasets(pool)

    file_path = os.path.expanduser("./datasets/ratings.csv")
    reader = Reader(line_format="user item rating", sep=",", skip_lines=1)

    data = Dataset.load_from_file(file_path, reader=reader)

    trainset = data.build_full_trainset()
    algo.fit(trainset)

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

def get_resource_nearest_neighbors(algo, resource_url):
    id_to_url, name_to_url = read_resources_urls()

    raw_id = name_to_url[resource_url]
    inner_id = algo.trainset.to_inner_iid(raw_id)

    neighbors = algo.get_neighbors(inner_id, k=1)
    neighbors = (
        algo.trainset.to_raw_iid(inner_id) for inner_id in neighbors
    )
    neighbors = (id_to_url[rid] for rid in neighbors)

    return neighbors

class HTTPRequestHandler(BaseHTTPRequestHandler):
    def __init__(self, pool, algo, *args, **kwargs):
        self.pool = pool
        self.algo = algo
        super().__init__(*args, **kwargs)

    def do_GET(self):
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        json_data = json.loads(post_data.decode('utf-8'))

        resource_url = json_data.get('resource', '')

        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        
        try:
            nearest_neighbors = get_resource_nearest_neighbors(algo, resource_url)
            response_json = json.dumps(nearest_neighbors)
            self.wfile.write(response_json.encode('utf-8'))
        except:
            response_json = json.dumps({ 'error': 'Cannot find any recommendations' })
            self.wfile.write(response_json.encode('utf-8'))

def run_server(pool, algo, port=8000):
    server_address = ('', port)
    httpd = HTTPServer(server_address, lambda *args, **kwargs: HTTPRequestHandler(pool, algo, *args, **kwargs))
    print(f'Starting server on port {port}...')
    httpd.serve_forever()

pool = create_database_pool()

sim_options = {"name": "cosine", "user_based": False}
algo = KNNBaseline(sim_options=sim_options)
train_knn(algo, pool)

run_server(pool, algo)