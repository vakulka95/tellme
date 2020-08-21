// Show and hide modal
const filter_btn = document.querySelector('.filter .filter-btn');
const close_btn = document.querySelector('.close');
const modal = document.querySelector('.modal');
const body = document.querySelector('body');
const user = document.querySelector('.user');
const userImg = document.querySelector('.user img');
const menuModal = document.querySelector('.menu-modal');
const desktopNav = document.querySelector('.desktop-nav')
const menuLink= document.querySelectorAll('.menu--link');
const navLink = document.querySelectorAll('.nav-link');
const tablesTr = document.querySelectorAll('table tr');
const filterSearch = document.querySelector('.search');
const nameInput = document.querySelectorAll('.name');
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
const dateOfReg = document.querySelector('.date-of-reg')
const description = document.querySelectorAll('.req-description');
const mobileDescription = document.querySelectorAll('.mobile-req-description');
const table = document.querySelector('.table');
const forMobile = document.querySelector('.for-mobile');
const aboutExpert = document.querySelector('.about-expert');
const status = document.querySelector('.table-cell-status');
const countOfPages = document.querySelector('.all');


// Remove string with count of page
let href = window.location.href;

if(href.indexOf('/requisition/') > 0 || href.indexOf('/expert/') > 0 || href.indexOf('/review/') > 0){
    countOfPages.classList.add('hide');
}else{
    countOfPages.classList.remove('hide');
}


if(aboutExpert){
    if(status.classList.contains('created')){
        aboutExpert.style.pointerEvents = 'none';
        aboutExpert.style.background = 'rgba(55, 169, 250, 0.4)';
    }else{
        aboutExpert.style.pointerEvents = 'auto';
        aboutExpert.style.background = 'rgba(55, 169, 250)';
    }
}

// Cut description
    if(table){
        description.forEach(el => {
            if(el.innerHTML.length > 45){
               let cut = el.innerText.substr(0, 45) + '...';
               el.innerText = cut;
            }
        })
    }
 
// Active link

for(let i = 0; i < navLink.length; i++){
    if(href.indexOf('/requisition') > 0){
        navLink[0].classList.add('nav-link--active');
    }
    else if(href.indexOf('/expert?') > 0 || href.indexOf('/expert/') > 0){
        navLink[1].classList.add('nav-link--active');
    }
    else if(href.indexOf('/profile') > 0){
        navLink[1].classList.add('nav-link--active');
        dateOfReg.style.width = '100%';
    }
    else if(href.indexOf('/review') > 0){
        navLink[2].classList.add('nav-link--active');
    }
    else if(href.indexOf('/expert_rating') > 0){
        navLink[3].classList.add('nav-link--active');
    }
}



// Edit block in expert item
if(edit){
        if(!editName){
            editSpec.style.width = '100%';
        }else{
            editSpec.style.width = '47%';
        }


        if(!editName && document.clientWidth <= 1024){
            modalContent.style.height = '50%';
        }
    edit.onclick = function(){
        modalEdit.style.display = 'flex';
        body.style.overflowY = 'hidden';
       
    }
        
}

if(close){
    close.onclick = function (){
        modalEdit.style.display = 'none';
        body.style.overflowY = 'scroll';
    }
}



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
            tablesTr[i].style.backgroundColor = '#F2F6FF';
        }else{
            tablesTr[i].style.backgroundColor = '#fff';
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
        }else{
            nameInput.forEach(el => {
                el.parentElement.parentElement.classList.remove('hide')
            });
        }
    }
}


  
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

