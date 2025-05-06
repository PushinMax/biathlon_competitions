package repository


import "fmt"


func (repo *Repo) StartRange(id int, firingRange int) error {
    comp, ok := repo.table[id]
    if !ok {
        return fmt.Errorf("Competitor(%d) is not registered", id)
    }
    if comp.StatusFlag != 1 {
        return fmt.Errorf("Competitor(%d) at another stage", id)    
    }

    comp.StatusFlag = 2
    comp.Goals = make([]bool, 5) 

    return nil
} 

func (repo *Repo) Hit(id, hit int) error {
    comp, ok := repo.table[id]
    if !ok {
        return fmt.Errorf("Competitor(%d) is not registered", id)
    }
    if comp.StatusFlag != 2 {
        return fmt.Errorf("Competitor(%d) at another stage", id)    
    }

    if comp.Goals[hit - 1] {
        return fmt.Errorf("Competitor(%d) already hit it(%d)", id, hit)    
    }

    comp.Goals[hit - 1] = true
    return nil
}

func (repo *Repo) LeftRange(id int) error {
    comp, ok := repo.table[id]
    if !ok {
        return fmt.Errorf("Competitor(%d) is not registered", id)
    }
    if comp.StatusFlag != 2 {
        return fmt.Errorf("Competitor(%d) at another stage", id)    
    }
    count := 0
    for _, v := range comp.Goals {
        if v {
            count++
        }
    }
    comp.Shots += 5
    comp.Hits += count

    comp.CurrentPenaltyCount = 5 - count
    if comp.CurrentPenaltyCount == 0 {
        comp.StatusFlag = 1
        return nil
    }


    comp.StatusFlag = 3
    return nil
}
