if (localStorage.getItem('token')) {
    // Redirect to dashboard if already logged in
    window.location.href = '/'; // Change to your actual dashboard URL
}
// Toggle password visibility
document.getElementById('togglePassword').addEventListener('click', function () {
    const passwordInput = document.getElementById('password');
    const icon = this.querySelector('i');

    if (passwordInput.type === 'password') {
        passwordInput.type = 'text';
        icon.classList.replace('fa-eye', 'fa-eye-slash');
    } else {
        passwordInput.type = 'password';
        icon.classList.replace('fa-eye-slash', 'fa-eye');
    }
});

// Form submission
document.getElementById('loginForm').addEventListener('submit', function (e) {
    e.preventDefault();

    // Show loading spinner
    document.getElementById('loadingSpinner').classList.remove('hidden');

    // Simulate authentication (replace with actual API call)
    setTimeout(() => {
        // Hide loading spinner
        document.getElementById('loadingSpinner').classList.add('hidden');

        // Get form values
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        // Simple validation (in a real app, you'd have proper validation)
        if (email && password) {
            fetch('/admin/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
            })
                .then(response => response.json())
                .then(data => {
                    console.log(data); // Log the response for debugging
                    if (data.message === 'Login successful!') {
                        localStorage.setItem('token', data.token); // Store token in local storage
                        // Redirect to dashboard on successful login
                        window.location.href = '/admin'; // Change to your actual dashboard URL
                    } else {
                        alert('Invalid email or password');
                    }
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert('An error occurred. Please try again later.');
                });
        } else {
            alert('Please enter both email and password');
        }
    }, 1500);
});

// Add floating animation to bus icon
const busIcon = document.querySelector('.bus-icon');
busIcon.addEventListener('mouseenter', () => {
    busIcon.style.animationPlayState = 'paused';
});
busIcon.addEventListener('mouseleave', () => {
    busIcon.style.animationPlayState = 'running';
});