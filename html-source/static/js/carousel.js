
// Get the carousel element
var carousel = document.getElementById('main-carousel');

// Create a function to start the carousel auto-play
function startCarousel() {
// Get all carousel items
var carouselItems = carousel.querySelectorAll('.carousel-item');

// Set the desired interval in milliseconds (e.g., 2000 for 2 seconds)
var interval = 3500;

// Loop through each carousel item
carouselItems.forEach(function(item, index) {
    // Set the data-bs-interval attribute of each item
    item.setAttribute('data-bs-interval', interval);
    
    // Set the first item as active
    if (index === 0) {
    item.classList.add('active');
    } else {
    item.classList.remove('active');
    }
});

// Initialize the carousel using the Bootstrap Carousel API
var carouselInstance = new bootstrap.Carousel(carousel, {
    interval: interval, // Set the desired interval in milliseconds
    pause: false,
    wrap: true // Allow wrapping from the last slide to the first slide (optional)
});
}

// Call the startCarousel function to begin auto-play
startCarousel();

