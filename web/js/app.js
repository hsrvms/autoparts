document.addEventListener("DOMContentLoaded", function () {
  console.log("AutoParts Inventory Management System loaded.");

  document.body.addEventListener("htmx:configRequest", function (event) {
    // Add any request configuration here
  });

  document.body.addEventListener("htmx:afterSwap", function (event) {
    // Code to run after HTMX content swap
  });

  document.body.addEventListener("htmx:responseError", function (event) {
    console.error("HTMX request error:", event.detail.error);
    // Handle errors, perhaps show a notification
  });
});

function formatCurrency(amount) {
  return new Intl.NumberFormat("tr-TR", {
    style: "currency",
    currency: "TRY",
  }).format(amount);
}

function formatDate(dateString) {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat("tr-TR", {
    year: "numeric",
    month: "short",
    day: "numeric",
  }).format(date);
}
