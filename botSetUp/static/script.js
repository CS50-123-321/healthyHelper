document.getElementById('habitForm').addEventListener('submit', async function(event) {
    event.preventDefault();

    const formData = {
        name: document.getElementById('name').value,
        habitName: document.getElementById('habitName').value,
        commitmentPeriod: document.getElementById('commitmentPeriod').value
    };

    try {
        const response = await fetch('/create-habit', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        const successMessage = document.getElementById('successMessage');

        if (response.ok) {
            updateProgressBar(100); // Animates to full after successful form submission

            // Show the success message
            successMessage.classList.remove('hide');
            successMessage.classList.add('show');

            setTimeout(() => {
                successMessage.classList.remove('show');
                successMessage.classList.add('hide');
            }, 3000); // Hide after 3 seconds
        } else {
            // Optional: handle error case if needed
            console.error('Error saving habit.');
        }
    } catch (error) {
        console.error('Error:', error);
    }
});

function updateProgressBar(percentage) {
    const progressBar = document.getElementById('progressBar');
    progressBar.style.width = `${percentage}%`;
}
