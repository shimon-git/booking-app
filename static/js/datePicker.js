// DatePicker function by id
const IDElem = document.getElementById('reservation-date');
const rangepicker = new DateRangePicker(IDElem, {
    // ...options
    format: "dd-mm-yyyy",
    minDate: new Date(),
}); 

