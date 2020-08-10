
// Show and hide modal
const filter_btn = document.querySelector('.filter .filter-btn');
const close_btn = document.querySelector('.close');
const modal = document.querySelector('.modal');
const body = document.querySelector('body');
const user = document.querySelector('.user');
const userImg = document.querySelector('.user img');
const menuModal = document.querySelector('.menu-modal');
const desktopNav = document.querySelector('.desktop-nav')
// const drag = document.querySelector('.drag');
// const dragTitle = document.querySelector('.drag-title');
// const filter = document.querySelector('.filter');
// const filterDesk = document.querySelector('.filter-desktop');
// const appClients = document.querySelector('.app__clients');
// const psychologists = document.querySelector('.psychologists');
// const header = document.querySelector('.header');
// const main = document.querySelector('.main');
// const rate = document.querySelector('.rate');
// const feedbacks = document.querySelector('.feedbacks');
const menuLink= document.querySelectorAll('.menu--link');
const tablesTr = document.querySelectorAll('table tr');
// const titleArrow = document.querySelector('.title__arrow');
// const button = document.querySelector('.button');
// const gender = document.querySelector('.app-gender-form');
// const form = document.querySelector('.form');
// const arrow = document.querySelector('.arrow');
const filterSearch = document.querySelector('.filter-desktop-search');
const nameInput = document.querySelectorAll('.name');
// const phoneInput = document.querySelectorAll('.phone');
const filterMobSearch = document.querySelector('.filter-mob-input');
const nameMob = document.querySelectorAll('.name-mob');
const eyeClose = document.querySelector('.eye-close');
const eyeOpen = document.querySelector('.eye-open');
const signInPass = document.getElementById('password');
const logoBlock = document.querySelector('.logo-row');
const modalEdit = document.querySelector('.modal-edit');
const modalContent = document.querySelector('.modal-content')
const edit = document.querySelector('.edit-btn');
const close = document.querySelector('.close');
const editName = document.querySelector('.edit-name');
const editSpec = document.querySelector('.edit-spec');

// Edit block in expert item
if(edit){
        if(!editName){
            editSpec.style.width = '100%';
        }else{
            editSpec.style.width = '47%';
        }
    edit.onclick = function(){
        modalEdit.style.display = 'flex';
        body.style.overflowY = 'hidden';
       
    }
        if(!editName && document.clientWidth <= 1024){
            modalContent.style.height = '50%';
        }
}

if(close){
    // document.onclick = function(){
    //     modalEdit.style.display = 'none';
    //     body.style.overflowY = 'scroll';
    // }
    close.onclick = function (){
        modalEdit.style.display = 'none';
        body.style.overflowY = 'scroll';
    }
}

// if(!clearable){
//     editSpec.style.width = '100%';
// }else{
//     editSpec.style.width = '47%';
// }


// Logo in log-in page
if(window.location.href == 'https://staging.tellme.com.ua/admin/login'){
    logoBlock.classList.remove('col-6');
    logoBlock.classList.add('col-12');
    logoBlock.classList.add('d-flex');
    logoBlock.classList.add('justify-content-center');
}else{
    logoBlock.classList.remove('col-12');
    logoBlock.classList.remove('d-flex');
    logoBlock.classList.remove('justify-content-center');
    logoBlock.classList.add('col-6');

}

// show modal menu
if(user){
    user.onclick = () => {
        menuModal.classList.toggle('show');
        menuModal.onclick = () => { 
        menuModal.classList.remove('show');
        }
    } 
}


if(filter_btn){
    filter_btn.onclick = showModal;

    function showModal(){
        modal.classList.add('show');
        body.style.overflowY = 'hidden';
    }
    
    close_btn.onclick = closeModal;
    
    function closeModal(){
        modal.classList.remove('show');
        body.style.overflowY = 'scroll';
    }    
}

// Active nav link

if(menuLink){
    menuLink.forEach(el => {
    el.onclick = () => {
        if(!el.classList.contains('menu--link--active')){
            menuLink.forEach(link => link.classList.remove('menu--link--active'));
            el.classList.add('menu--link--active');
        }
    }
})}


// Background table


function spillTr(){
    for(let i = 1; i < tablesTr.length; i++){
        if(i % 2 === 0 ){
            tablesTr[i].style.backgroundColor = '#E3E8F5';
        }
    }
}

spillTr();

// SEARCH NAME 
if(filterSearch){
    filterSearch.oninput = function(){
        let val = this.value.trim();
        console.log(val)
    
        if(val != ''){
            nameInput.forEach(el => {
                if(el.innerText.search(val) == -1 ){
                    el.parentElement.parentElement.classList.add('hide')
                }
                else{
                    el.parentElement.parentElement.classList.remove('hide');
                }
            });
            // phoneInput.forEach(el => {
            //     if(el.innerText.search(val) == -1 ){
            //         el.parentElement.parentElement.classList.add('hide')
            //     }
            //     else{
            //         el.parentElement.parentElement.classList.remove('hide');
            //     }
            // });
        }else{
            nameInput.forEach(el => {
                el.parentElement.parentElement.classList.remove('hide')
            });
            // phoneInput.forEach(el => {
            //     el.parentElement.parentElement.classList.remove('hide')
            // });
        }
    }
}
// filterSearch.oninput = function(){
//     let val = this.value.trim();
//     console.log(val)

//     if(val != ''){
//         nameInput.forEach(el => {
//             if(el.innerText.search(val) == -1 ){
//                 console.log( el.parentElement.parentElement);
//                 el.parentElement.parentElement.parentElement.classList.add('hide')
//             }
//             else{
//                 el.parentElement.parentElement.parentElement.classList.remove('hide');
//             }
//         });
    
//     }else{
//         nameInput.forEach(el => {
//             el.parentElement.parentElement.parentElement.classList.remove('hide')
//         });
        
//     }
// }


// Search application in search input


  
if(filterMobSearch){
    filterMobSearch.oninput = function(){
        let val = this.value.trim();
        if(val != ''){
            nameMob.forEach(el => {
                if(el.innerText.search(val) == -1 ){
                    el.parentElement.parentElement.parentElement.parentElement.parentElement.classList.add('hide')
                }
                else{
                    el.parentElement.parentElement.parentElement.parentElement.parentElement.classList.remove('hide');
                }
            });
            
        }else{
            nameMob.forEach(el => {
                el.parentElement.parentElement.parentElement.parentElement.parentElement.classList.remove('hide')
            });
            
        }
    }
}





// $(document).ready(function($) {
//     $(".clickable-row").click(function() {
//         window.document.location = $(this).data("href");
//     });
// });

if(eyeClose){
    eyeClose.onclick = function(){
        if( !eyeClose.classList.contains('hide')){
            eyeClose.classList.add('hide');
            eyeOpen.classList.remove('hide')
            signInPass.setAttribute('type', 'text');
        }
        eyeOpen.onclick = function(){
            eyeClose.classList.remove('hide');
            eyeOpen.classList.add('hide')
            signInPass.setAttribute('type', 'password');

        }
    }

}

// function onSubmitLogin(token) {
//     document.getElementById('reCaptchaForm').submit();
// }







    


// change Arrow

// function changeArrow(){
//     if(document.body.clientWidth >= 1024){
//         form.style.width = '75%';
//         userImg.style.display = 'none';
//         user.style.display = 'none';
//         desktopNav.style.display = 'flex';
//         gender.style.display = 'none';

//     }else{
//         form.style.width = '100%';
//         userImg.style.display = 'block';
//         user.style.display = 'block';
//         desktopNav.style.display = 'none';
//         gender.style.display = 'block';
//     }
// }
// changeArrow();

// window.addEventListener('resize', changeArrow); 






// function hideBlock(){
//     if(document.body.clientWidth >= 1024){
//         drag.classList.add('d-flex');
//         dragTitle.style.display = 'block';
//         user.style.display = 'none';
//         userImg.style.display = 'none';
//         desktopNav.style.display = 'block';
//     }else{
//         drag.classList.remove('d-flex')
//         drag.style.display = 'none';
//         dragTitle.style.display = 'none';
//         user.style.display = 'block';
//         userImg.style.display = 'block';
//         desktopNav.style.display = 'none';

//     }
// }
// hideBlock();

// window.addEventListener('resize', hideBlock);


// Edit inputs


// Hide elements



// function hideFilterBlock(){
//     if(document.body.clientWidth >= 1024){
//         filter.style.display = 'none';
//         filterDesk.style.display = 'block';
//         appClients.style.display = 'none';
//         header.style.position = 'relative';
//         main.style.display = 'block';
//         userImg.style.display = 'none';
//         user.style.display = 'none';
//         desktopNav.style.display = 'block';
//         psychologists.style.display = 'none';
//         feedbacks.style.display = 'none';

//     }else{
//         filter.style.display = 'block';
//         filterDesk.style.display = 'none';
//         appClients.style.display = 'block';
//         header.style.position = 'sticky';
//         main.style.display = 'none';
//         userImg.style.display = 'block';
//         desktopNav.style.display = 'none';
//         psychologists.style.display = 'block';
//         feedbacks.style.display = 'block';
//     }
// }
// hideFilterBlock();

// window.addEventListener('resize', hideFilterBlock);
