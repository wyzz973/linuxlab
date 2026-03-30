#!/bin/bash
expected='# DONE: Add error handling
def connect_db(host, port):
    conn = create_connection(host, port)
    if conn is None:
        raise ConnectionError("Failed")
    return conn

# DONE: Add retry logic
def fetch_data(conn, query):
    cursor = conn.execute(query)
    return cursor.fetchall()'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
