from flask import Flask, request, jsonify
from supabase import create_client, Client
import os

app = Flask(__name__)

SUPABASE_URL = "https://utzdhilbitfcdeljnctj.supabase.co"
SUPABASE_KEY = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InV0emRoaWxiaXRmY2RlbGpuY3RqIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MTE3MjAyOTQsImV4cCI6MjAyNzI5NjI5NH0.oTHeYvKzSEGcoBu8pbMMeDWvbzmwiEFIzEUQBUmbgKk"

db: Client = create_client(SUPABASE_URL, SUPABASE_KEY)

@app.route("/", methods=["POST"])
def add_data():
    data = request.json
    name = data.get("name")
    email = data.get("email")
    phone_number = data.get("phone_number")
    address = data.get("address")

    response = db.table("api_uts").insert({"name": name, "email": email, "phone_number": phone_number, "address": address}).execute()
    return jsonify(response.data)

@app.route("/", methods=["GET"])
def get_data():
    response = db.table("api_uts").select("*").execute()

    return jsonify(response.data)


if __name__ == "__main__":
    app.run(port=3000, debug=True)
