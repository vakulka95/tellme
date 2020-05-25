from locust import HttpLocust, TaskSet

def reg(l):
    payload = {"name": "тест запит", "gender": "male", "age": 20, "phone": "+380638888888", "diagnosis": "panic",
     "diagnosisDescription": "омагад-омагад-омагад-омагад-омагад-омагад-омагад-омагад-омагад-омагад-омагад-омагад-омагад-омагад-омагад-омагад-омагад-омагад",
     "expertGender": "female", "feedbackType": "phone", "feedbackTime": "8:00", "feedbackWeekDay": "mon", "isAdult": True}
    l.client.post("/api/v1/requisition", json=payload)

class UserBehavior(TaskSet): 
    tasks = {reg: 1}

class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 5000
    max_wait = 9000