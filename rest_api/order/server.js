const express = require('express');
const axios = require('axios');

const app = express();
const PORT = 3001;

const SERVICE_PORT = 3000;

app.get('/', async (req, res) => {
    const { data } = await axios.get(`http://localhost:${SERVICE_PORT}`);

    res.json(data);
});

app.listen(PORT, () => {
  console.log(`server started on port ${PORT}`);
});