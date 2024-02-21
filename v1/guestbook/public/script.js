$(document).ready(function() {
  // Selecting necessary elements from the DOM
  var headerTitleElement = $("#header h1");
  var entriesElement = $("#guestbook-entries");
  var formElement = $("#guestbook-form");
  var submitElement = $("#guestbook-submit");
  var entryContentElement = $("#guestbook-entry-content");
  var hostAddressElement = $("#guestbook-host-address");

  // Function to append guestbook entries to the DOM
  var appendGuestbookEntries = function(data) {
    entriesElement.empty(); // Clearing existing entries
    $.each(data, function(key, val) {
      entriesElement.append("<p>" + val + "</p>"); // Appending new entries
    });
  }

  // Function to handle form submission
  var handleSubmission = function(e) {
    e.preventDefault(); // Preventing default form submission
    var entryValue = entryContentElement.val(); // Getting the value from the input field
    if (entryValue.length > 0) {
      entriesElement.append("<p>...</p>"); // Adding a temporary entry indicator
      $.getJSON("rpush/guestbook/" + entryValue, appendGuestbookEntries); // Sending data to the server and updating guestbook
      entryContentElement.val(""); // Clearing the input field after submission
    }
    return false;
  }

  // Click event handler for the submit button
  submitElement.click(handleSubmission);

  // Submit event handler for the form (in case Enter is pressed)
  formElement.submit(handleSubmission);

  // Displaying the host address in the designated element
  hostAddressElement.append(document.URL);

  // Polling the server every second to update guestbook entries
  (function fetchGuestbook() {
    $.getJSON("lrange/guestbook").done(appendGuestbookEntries).always(
      function() {
        setTimeout(fetchGuestbook, 1000); // Polling interval of 1 second
      });
  })();
});
