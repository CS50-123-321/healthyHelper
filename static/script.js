// Extract the session ID from the URL
function getSessionID() {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get('session');
}

const sessionID = getSessionID(); 

document.getElementById('habitForm').addEventListener('submit', async function (event) {
    event.preventDefault();

    const formData = {
        name: document.getElementById('name').value,
        habit_name: document.getElementById('habit_name').value,
        commitment_Period: document.getElementById('commitment_Period').value,
        tele_id: sessionID,
        solo: document.getElementById('solo').value

    };
    try {
        const response = await fetch('/create-habit', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json', // Ensure this matches what your server expects
            },
            body: JSON.stringify(formData)
        });

        // Handle the server response
        if (response.ok) {
            const successMessage = document.getElementById('successMessage');
            updateProgressBar(100); // Animates to full after successful form submission

            // Show the success message
            successMessage.classList.remove('hide');
            successMessage.classList.add('show');

            setTimeout(() => {
                successMessage.classList.remove('show');
                successMessage.classList.add('hide');
            }, 3000); // Hide after 3 seconds
        } else {
            // Handle error case if needed
            console.error('Error saving habit.');
        }
    } catch (error) {
        console.error('Error:', error);
    }
});

// Function to update the progress bar (can be used after form submission)
function updateProgressBar(percentage) {
    const progressBar = document.getElementById('progressBar');
    progressBar.style.width = `${percentage}%`;
}
