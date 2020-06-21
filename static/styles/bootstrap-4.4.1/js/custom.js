function isDo(text){
    return confirm(text);
}

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
    });
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

    const ok = isDo("Ви впевнені що хочете видалити психолога?");

    if(ok) {
        $.ajax({
            url: '/admin/expert/' + expertId,
            type: 'DELETE',
            success: function (result) {
                history.back();
            },
            error: function (error) {
                console.log(error);
            }
        })
    }
}

function deleteExpertDocumentAPI(expert){
    let expertId = $(expert).data('expert-id');
    let documentId = $(expert).data('document-id');

    const ok = isDo("Ви впевнені що хочете видалити документ?");

    if(ok) {
        $.ajax({
            url: '/admin/expert/' + expertId + '/document/' + documentId,
            type: 'DELETE',
            success: function (result) {
                history.back();
            },
            error: function (error) {
                console.log(error);
            }
        })
    }
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

function deleteRequisitionAPI(requisition){
    let requisitionId = $(requisition).data('requisition-id');

    const ok = isDo("Ви впевнені що хочете видалити заявку?");

    if(ok) {
        $.ajax({
            url: '/admin/requisition/' + requisitionId,
            type: 'DELETE',
            success: function (result) {
                history.back();
            },
            error: function (error) {
                console.log(error);
            }
        })
    }
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


function discardRequisitionAPI(requisition){
    let requisitionId = $(requisition).data('requisition-id');

    $.ajax({
        url: '/admin/requisition/'+requisitionId+'/discard',
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

function createSessionRequisitionAPI(requisition){
    let requisitionId = $(requisition).data('requisition-id');

    $.ajax({
        url: '/admin/requisition/'+requisitionId+'/session/apply',
        type: 'POST',
        success: function (result) {
            document.location.reload();
        },
        error: function (error) {
            console.log(error);
        }
    })
}

function discardSessionRequisitionAPI(requisition){
    let requisitionId = $(requisition).data('requisition-id');

    $.ajax({
        url: '/admin/requisition/'+requisitionId+'/session/discard',
        type: 'DELETE',
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

function tog(v){return v?'addClass':'removeClass';}
$(document).on('input', '.clearable', function(){
    $(this)[tog(this.value)]('x');
}).on('mousemove', '.x', function( e ){
    $(this)[tog(this.offsetWidth-18 < e.clientX-this.getBoundingClientRect().left)]('onX');
}).on('click', '.onX', function(){
    $(this).removeClass('x onX').val('').change();
});

function onSubmitLogin(token) {
    document.getElementById('reCaptchaForm').submit();
}
