package main
import "sync"

type Blackboard struct {
    board map[string][]interface{}
    L sync.Mutex;
    C *sync.Cond
}

func NewBlackboard() *Blackboard {
    var b Blackboard;
    b.board = make(map[string][]interface{});
    b.C = sync.NewCond(&b.L)
    return &b;
}

func (self *Blackboard) get(tags ...string) (string, interface{}) {
    var tag string
    
    self.L.Lock()
    defer self.L.Unlock()
    
    found := false
    for !found {
        for _,tag = range tags {
            _,present := self.board[tag]
            if present {
                found = true
                break
            }
        }
        if !found {
            self.C.Wait()
        }
    }
    
    tmp := self.board[tag][0]
    if len(self.board[tag]) == 1 {
        delete(self.board,tag)
    } else {
        self.board[tag] = self.board[tag][1:]
    }
    return tag,tmp
}


func (self *Blackboard) put(tag string, data interface{}) {
    self.L.Lock()
    defer self.L.Unlock()
    _,present := self.board[tag]
    if !present {
        self.board[tag] = make([]interface{}, 1)
        self.board[tag][0] = data
    } else {
        self.board[tag] = append( self.board[tag], data )
    }
    self.C.Broadcast()
}


