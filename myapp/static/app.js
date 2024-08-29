document.addEventListener('DOMContentLoaded', function() {
    // Index Page Elements
    const loginForm = document.getElementById('loginForm');
    const signupForm = document.getElementById('signupForm');
    const toggleFormButton = document.getElementById('toggleFormButton');
    const formTitle = document.getElementById('formTitle');
    const message = document.getElementById('message');

    // Congrats Page Elements
    const logoutButton = document.getElementById('logoutButton');

    // If on the index.html page
    if (toggleFormButton) {
        // Handle form toggle
        toggleFormButton.addEventListener('click', function() {
            if (signupForm && signupForm.style.display === 'none') {
                // Switch to Sign-Up
                signupForm.style.display = 'block';
                if (loginForm) loginForm.style.display = 'none';
                if (formTitle) formTitle.textContent = 'Sign Up';
                toggleFormButton.textContent = 'Already have an account? Log In';
            } else {
                // Switch to Log-In
                if (signupForm) signupForm.style.display = 'none';
                if (loginForm) loginForm.style.display = 'block';
                if (formTitle) formTitle.textContent = 'Log In';
                toggleFormButton.textContent = "Don't have an account? Sign Up";
            }
        });

        if (signupForm) {
            // Handle Sign-Up form submission
            signupForm.addEventListener('submit', function(e) {
                e.preventDefault();
                const email = document.getElementById('signupEmail').value;
                const password = document.getElementById('signupPassword').value;

                fetch('/signup', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ email: email, password: password })
                })
                .then(response => {
                    if (!response.ok) {
                        return response.json().then(data => {
                            throw new Error(data.message);
                        });
                    }
                    return response.json();
                })
                .then(data => {
                    if (message) {
                        message.textContent = data.message;
                    }
                })
                .catch(error => {
                    if (message) {
                        message.textContent = error.message;
                    }
                    console.error("Sign up error: ", error);
                });
            });
        }

        if (loginForm) {
            // Handle Log-In form submission
            loginForm.addEventListener('submit', function(e) {
                e.preventDefault();
                const email = document.getElementById('loginEmail').value;
                const password = document.getElementById('loginPassword').value;

                fetch('/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ email: email, password: password })
                })
                .then(response => {
                    if (response.redirected) {
                        window.location.href = response.url;
                    }
                    else if (!response.ok){
                        return response.text().then(text => { throw new Error(text); });
                    }
                    else {
                        return response.json();
                    }
                })
                .then(data => {
                    if (data && data.message) {
                        if (message) {
                            message.textContent = data.message;
                        }
                    }
                })
                .catch(error => {
                    if (message) {
                        message.textContent = error.message;
                    }
                });
            });
        }
    }

    // If on the congrats.html page
    if (logoutButton) {
        logoutButton.addEventListener('click', function() {
            fetch('/logout', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
            })
            .then(response => {
                if (response.redirected) {
                    window.location.href = response.url;
                } else {
                    return response.json();
                }
            })
            .catch(error => {
                console.error('Error during logout:', error);
            });
        });
    }
});
