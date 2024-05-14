const express = require("express");
const bodyParser = require("body-parser");
const app = express();
const Producer = require("./producer");
const producer = new Producer();
const PORT = 3000;
const supabase = require("@supabase/supabase-js");
app.use(bodyParser.json("application/json"));

const Supabase_URl = "https://utzdhilbitfcdeljnctj.supabase.co";
const Supabase_Key = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InV0emRoaWxiaXRmY2RlbGpuY3RqIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MTE3MjAyOTQsImV4cCI6MjAyNzI5NjI5NH0.oTHeYvKzSEGcoBu8pbMMeDWvbzmwiEFIzEUQBUmbgKk";
const db = supabase.createClient(Supabase_URl, Supabase_Key);

app.post("/", async (req, res) => {
  const { name, email, phone_number, address } = req.body;
  const Post = await db.from("api_uts").insert({ name, email, phone_number, address });
  res.json({ 
    Post
  });
  const { data: allData } = await db
  .from('api_uts')
  .select('*');

  for (const mesage of allData) {
    await producer.publishMessage( mesage);
  }
  res.send();

});

app.listen(PORT, () => {
  console.log(` Server started on port ${PORT} `);
});
