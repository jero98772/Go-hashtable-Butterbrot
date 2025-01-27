document.getElementById('putForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const key = document.getElementById('putKey').value;
    const value = document.getElementById('putValue').value;

    const response = await fetch('/api/put', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ [key]: value }),
    });

    if (response.ok) {
        alert('Data inserted successfully');
    } else {
        alert('Error inserting data');
    }
});

document.getElementById('getForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const key = document.getElementById('getKey').value;

    const response = await fetch(`/api/get/${key}`);
    const result = await response.json();

    if (response.ok) {
        document.getElementById('getResult').innerText = JSON.stringify(result, null, 2);
    } else {
        document.getElementById('getResult').innerText = 'Key not found';
    }
});

document.getElementById('deleteForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const key = document.getElementById('deleteKey').value;

    const response = await fetch(`/api/delete/${key}`, {
        method: 'DELETE',
    });

    if (response.ok) {
        alert('Key deleted successfully');
    } else {
        alert('Error deleting key');
    }
});


document.getElementById('refreshData').addEventListener('click', async () => {
    const response = await fetch('/api/elements');
    if (response.ok) {
        const data = await response.json();
        document.getElementById('dhtElements').innerText = JSON.stringify(data.dht, null, 2);
        document.getElementById('redisElements').innerText = JSON.stringify(data.redis, null, 2);
    } else {
        alert('Error fetching data');
    }
});
