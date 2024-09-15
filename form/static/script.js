document.getElementById('habitForm').addEventListener('submit', async function(event) {
    event.preventDefault();

    const formData = {
        name: document.getElementById('name').value,
        habitName: document.getElementById('habitName').value,
        commitmentPeriod: document.getElementById('commitmentPeriod').value
    };

    try {
        const response = await fetch('/save-habit', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            alert('Habit saved successfully!');
        } else {
            alert('Error saving habit.');
        }
    } catch (error) {
        console.error('Error:', error);
    }
});
