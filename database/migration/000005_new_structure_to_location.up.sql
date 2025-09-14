CREATE TABLE IF NOT EXISTS task (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    state INT NOT NULL,
    completion_date TIMESTAMP NULL
);
/*
0 - not done
1 - done
*/
CREATE TABLE IF NOT EXISTS user_to_task (
    user_id UUID REFERENCES  users(id) NOT NULL,
    task_id UUID REFERENCES task(id) NOT NULL,
    UNIQUE(task_id)
);

CREATE TABLE IF NOT EXISTS current_user_task (
    task_id UUID REFERENCES task(id) NOT NULL,
    user_id UUID REFERENCES users(id) NOT NULL,
    PRIMARY KEY (task_id,user_id)

)

/*
Musze tez sledzic dlugosc pracy znaczy przedzialy, np pracownik w apce kliknie koniec pracy
i wtedy np current task jest napisane e poza praca albo cos takiego

albo zrobic logera co 1minut ktory do las_time_online bedzie updatowal i wtedy jesli nie bedzie dlugo update
to samo sie zmieni
*/
/*
jesli koncze jakies zadanie w mobile app to znika z tej ostatniej tabeli ale
nie wmoze to wplynac na pinsy wsensie pin zosatje ale jego opisy znikaja moze
*/