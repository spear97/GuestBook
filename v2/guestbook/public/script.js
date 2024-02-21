$(document).ready(function() {
  // Selecting elements from the DOM
  var headerTitleElement = $("#header h1"); // Selects the h1 element inside #header
  var entriesElement = $("#guestbook-entries"); // Selects the div with id #guestbook-entries
  var formElement = $("#guestbook-form"); // Selects the form with id #guestbook-form
  var submitElement = $("#guestbook-submit"); // Selects the submit button with id #guestbook-submit
  var entryContentElement = $("#guestbook-entry-content"); // Selects the input field with id #guestbook-entry-content
  var hostAddressElement = $("#guestbook-host-address"); // Selects the h2 element with id #guestbook-host-address

  // Function to append guestbook entries to the DOM
  var appendGuestbookEntries = function(data) {
    entriesElement.empty(); // Clear existing entries
    // Iterate through the data and append each entry as a <p> element
    $.each(data, function(key, val) {
      entriesElement.append("<p>" + val + "</p>");
    });
  }

  // Function to handle form submission
  var handleSubmission = function(e) {
    e.preventDefault(); // Prevent default form submission
    var entryValue = entryContentElement.val(); // Get the value entered in the input field
    if (entryValue.length > 0) {
      entriesElement.append("<p>...</p>"); // Add a temporary loading indicator
      // Send an AJAX request to add the entry to the guestbook
      $.getJSON("rpush/guestbook/" + entryValue, appendGuestbookEntries);
      entryContentElement.val(""); // Clear the input field after submission
    }
    return false; // Prevent default form behavior
  }

  // Attach event listeners to the submit button and form
  submitElement.click(handleSubmission); // Handle click on submit button
  formElement.submit(handleSubmission); // Handle form submission

  // Display the current host address in the guestbook
  hostAddressElement.append(document.URL);

  // Function to fetch guestbook entries periodically
  (function fetchGuestbook() {
    // Send an AJAX request to get the guestbook entries
    $.getJSON("lrange/guestbook").done(appendGuestbookEntries).always(
      // After getting the entries, continue to fetch every second
      function() {
        setTimeout(fetchGuestbook, 1000); // Fetch again after 1 second
      });
  })();
});
