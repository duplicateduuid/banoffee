from setuptools import setup, find_packages

setup(
    name='resources-recommender',
    version='0.1',
    packages=find_packages(),
    install_requires=[
        'setuptools',
        'numpy',
        'wheel',
        'scikit-surprise',
        'psycopg2-binary',
        'types-psycopg2',
        'mypy'
    ],
)