
// function to addindg toggle functionality to the nav-bar

document.addEventListener("DOMContentLoaded", function() {
    const dropdownToggle = document.querySelector(".dropdown-toggle");
    const dropdownMenu = document.querySelector(".dropdown-menu");

    dropdownToggle.addEventListener("click", function() {
        dropdownMenu.classList.toggle("show");
    });

    document.addEventListener("click", function(event) {
        if (!dropdownToggle.contains(event.target)) {
            dropdownMenu.classList.remove("show");
        }
    });
});








// Functions for alerts

function notify(msg,msgType) {
    notie.alert({
        type: msgType,
        text: msg,
    })
};

function notifyModal(title,text,icon,confirmButtonText) {
    Swal.fire({
        title: title,
        html: text,
        icon: icon,
        confirmButtonText: confirmButtonText
    })
}




function Prompt() {
let toast = function(c){
    const {
        msg = "",
        icon = "success",
        position = "top-end",

    } = c;
    const Toast = Swal.mixin({
        toast: true,
        title: msg,
        position: position,
        icon: icon,
        showConfirmButton: false,
        timer: 2000,
        timerProgressBar: true,
        didOpen: (toast) => {
            toast.addEventListener('mouseenter', Swal.stopTimer)
            toast.addEventListener('mouseleave', Swal.resumeTimer)
        }
    })

    Toast.fire({ })
};
let success = function(c){
    const {
        msg = "",
        title = "",
        footer = ""
    } = c;
    Swal.fire({
        icon: 'success',
        title: title,
        text: msg,
        footer: footer,
    })
};
let error = function(c){
    const {
        msg = "",
        title = "",
        footer = ""
    } = c;
    Swal.fire({
        icon: 'error',
        title: title,
        text: msg,
        footer: footer,
    })
};
async function custome(c){
    const {
        msg = "",
        title = "",
    } = c;
    const { value: result } = await Swal.fire({
        title: title,
        html: msg,
        backdrop: false,
        focusConfirm: false,
        showCancelButton: true,
        willOpen: () => {
            const elem = document.getElementById("reservation-dates-modal");
            const rp = new DateRangePicker(elem,{
                format: 'dd-mm-yyyy',
                minDate: new Date(),
                showOnFocus: true,
                orientation: 'top',
                container: '.modal-dialog',
            })
        },
        didOpen: () => {
            document.getElementById('start').removeAttribute('disabled')
            document.getElementById('end').removeAttribute('disabled')
        },
        preConfirm: () => {
            return [
            document.getElementById('start').value,
            document.getElementById('end').value
            ]
        }
    })
    if (result) {
        if (result.dismiss !== Swal.DismissReason.cancel) {
            if (result.value !== "") {
                if (c.callback !== undefined) {
                    c.callback(result);
                }
            } else {
                c.callback(false)
            }
        } else {
            c.callback(false)
        }
    }
}; 
return{
    toast: toast,
    success: success,
    error: error,
    custome: custome,
}
}







