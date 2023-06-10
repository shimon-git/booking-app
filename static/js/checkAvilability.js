let attantion = Prompt();
document.getElementById("check-avilability-button").addEventListener("click", function(){

let html = `
<form action="" method="post" novalidate class="needs-validation" style="overflow-x: hidden;">
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
attantion.custome({msg:html, title: "Choose Your Dates"})
});