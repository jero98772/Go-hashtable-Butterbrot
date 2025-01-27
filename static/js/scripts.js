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
// Function to transform flat object into nested node data
function transformToNodeData(data) {
    const nodeData = {};
    
    // Iterate through all keys in the data
    Object.entries(data).forEach(([key, value]) => {
        // Split the key into nodeHash and actual key
        const [nodeHash, actualKey] = key.split(':');
        
        // If this nodeHash doesn't exist in nodeData yet, create it
        if (!nodeData[nodeHash]) {
            nodeData[nodeHash] = {};
        }
        
        // Add the key-value pair to the appropriate nodeHash object
        nodeData[nodeHash][actualKey] = value;
    });
    
    return nodeData;
}


async function updateDHTDisplay() {
    try {
        const response = await fetch('/api/elements');
        const data = await response.json();
        console.log(data,data["dht"]);
        // Initialize objects for each node
        const nodeData = transformToNodeData(data["dht"]);

        // Sort data into appropriate nodes
        Object.entries(data).forEach(([key, value]) => {
            // Check which node this key belongs to
            for (const nodeHash of Object.keys(nodeData)) {
                if (key.startsWith(nodeHash)) {
                    nodeData[nodeHash][key] = value;
                    break;
                }
            }
        });

        // Get the container
        const container = document.getElementById('dhtNodesContainer');
        container.innerHTML = '';

        // Create and append frames for each node
        Object.entries(nodeData).forEach(([nodeHash, elements], index) => {
            const nodeNumber = index + 1;
            const nodeFrame = document.createElement('div');
            nodeFrame.id = `node${nodeNumber}Display`;
            nodeFrame.className = 'p-4 bg-white rounded shadow';
            
            nodeFrame.innerHTML = `
                <h3 class="text-lg font-medium">DHT Node ${nodeNumber} Elements</h3>
                <pre id="node${nodeNumber}Elements" class="bg-gray-100 p-2 rounded">${
                    JSON.stringify(elements, null, 2)
                }</pre>
            `;
            
            container.appendChild(nodeFrame);
        });

    } catch (error) {
        console.error('Failed to fetch DHT elements:', error);
    }
}


// Call the function to populate the frames
updateDHTDisplay();
