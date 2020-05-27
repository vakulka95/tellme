function blockExpertAPI(expert){
    let expertId = $(expert).data('expert-id');

    $.ajax({
        url: '/admin/expert/'+expertId+'/block',
        type: 'PUT',
        success: function (result) {
            document.location.reload();
        },
        error: function (error) {
            console.log(error);
        }
    })
}


function activateExpertAPI(expert){
    let expertId = $(expert).data('expert-id');

    $.ajax({
        url: '/admin/expert/'+expertId+'/activate',
        type: 'PUT',
        success: function (result) {
            document.location.reload();
        },
        error: function (error) {
            console.log(error);
        }
    })
}

function deleteExpertAPI(expert){
    let expertId = $(expert).data('expert-id');

    $.ajax({
        url: '/admin/expert/'+expertId,
        type: 'DELETE',
        success: function (result) {
            history.back();
        },
        error: function (error) {
            console.log(error);
        }
    })
}

function changePasswordExpertAPI(expert){
    let expertId = $(expert).data('expert-id');
    let pwd = document.getElementById('pwd-' + expertId);

    $.ajax({
        url: '/admin/expert/'+expertId+'/password',
        type: 'PUT',
        data: `{"password": "${pwd.value}"}`,
        success: function (result) {
            let res = document.getElementById('changePasswordBlock-' + expertId);
            res.classList.add('alert-success');
            res.innerText = 'Пароль змінено';
        },
        error: function (error) {
            let res = document.getElementById('changePasswordBlock-' + expertId);
            res.classList.add('alert-danger');
            res.innerText = 'Помилка при зміні пароля';
        }
    })
}

function takeRequisitionAPI(requisition){
    let requisitionId = $(requisition).data('requisition-id');

    $.ajax({
        url: '/admin/requisition/'+requisitionId+'/take',
        type: 'PUT',
        success: function (result) {
            document.location.reload();
        },
        error: function (error) {
            console.log(error);
        }
    })
}

function completeRequisitionAPI(requisition){
    let requisitionId = $(requisition).data('requisition-id');

    $.ajax({
        url: '/admin/requisition/'+requisitionId+'/complete',
        type: 'PUT',
        success: function (result) {
            document.location.reload();
        },
        error: function (error) {
            console.log(error);
        }
    })
}

$(document).ready(function($) {
    $(".table-div").click(function() {
        window.document.location = $(this).data("href");
    });
});
