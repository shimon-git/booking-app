let attantion = Prompt();
document.getElementById("check-avilability-button").addEventListener("click", function(){

let html = `
<form id="check-avilability-form" action="" method="post" novalidate class="needs-validation" style="overflow-x: hidden;">
<div class="row" id="reservation-dates-modal">
    <div class="col" style="padding-right: 15px; padding-left: 15px;">
        <label for="start_date" class="form-label">Starting Date</label>
        <input type="text" class="form-control" id="start" name="start" disabled required autocomplete="off" placeholder="Arrival date" >
    </div>
    <div class="col" style="padding-right: 15px; padding-left: 15px;">
        <label for="end_date" class="form-label">Ending Date</label>
        <input type="text" class="form-control" id="end" name="end" disabled required autocomplete="off" placeholder="Departure">
    </div>
</div>
</form>
`;


// Get the CSRF token from the data attribute
var csrfToken = document.querySelector('script[src="/static/js/checkAvilability.js"]').getAttribute('data-csrf');

attantion.custome({
    msg:html,
    title: "Choose Your Dates",
    callback: function(result) {
         console.log("called");
         let form = document.getElementById('check-avilability-form');
         let formData = new FormData(form)
         formData.append("csrf_token", csrfToken);

         fetch('/check-avilability-json', {
            method: "post",
            body: formData
         })
         .then(response => response.json())
         .then(data => {
            console.log(data)
            console.log(data.ok)
            console.log(data.message)
         })
    }
})
});