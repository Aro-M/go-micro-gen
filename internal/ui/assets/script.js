document.getElementById('generator-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const statusEl = document.getElementById('status-message');
    const btn = document.getElementById('generate-btn');

    // UI Loading state
    btn.innerHTML = 'Generating... ⏳';
    btn.disabled = true;
    statusEl.className = 'status-message'; // reset
    statusEl.innerHTML = '';

    // Serialize Payload
    const payload = {
        name: document.getElementById('name').value.trim(),
        module: document.getElementById('module').value.trim(),
        db: document.getElementById('db').value,
        broker: document.getElementById('broker').value,
        cloud: document.getElementById('cloud').value,
        serverless: document.getElementById('serverless').checked,
        graphql: document.getElementById('graphql').checked,
        jwt: document.getElementById('jwt').checked,
        seeding: document.getElementById('seeding').checked
    };

    try {
        const res = await fetch('/api/generate', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        const data = await res.json();

        statusEl.classList.add('show');
        if (res.ok && data.status === 'success') {
            statusEl.classList.add('status-success');
            statusEl.innerText = `✅ ${data.message}`;
        } else {
            statusEl.classList.add('status-error');
            statusEl.innerText = `❌ Error: ${data.error || 'Failed to generate microservice.'}`;
        }
    } catch (err) {
        statusEl.classList.add('show');
        statusEl.classList.add('status-error');
        statusEl.innerText = `❌ Request Failed: ${err.message}`;
    } finally {
        btn.innerHTML = '✨ Generate Microservice';
        btn.disabled = false;
    }
});
