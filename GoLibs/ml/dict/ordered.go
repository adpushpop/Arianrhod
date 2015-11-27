package dict

import (
    . "fmt"
    . "ml/array"
)

type orderedDict struct {
    dict Dict
    keys Array
}

type orderedDictItem struct {
    Key     interface{}
    Value   interface{}
}

type OrderedDict struct {
    *orderedDict
}

func NewOrderedDict() *OrderedDict {
    return &OrderedDict{
        orderedDict : &orderedDict{
            dict: Dict{},
            keys: Array{},
        },
    }
}

func (self *orderedDict) toString(depth int) string {
    space := ""

    for i := depth; i > 0; i-- {
        space += "  "
    }

    s := "{\n"
    for _, k := range self.Keys() {
        var key, value string

        v := self.Get(k)

        switch t := k.(type) {
            case string:
                key = Sprintf("'%v'", t)

            default:
                key = Sprintf("%v", t)
        }

        switch obj := v.(type) {
            case *orderedDict:
                value = obj.toString(depth + 1)

            case orderedDict:
                value = obj.toString(depth + 1)

            case Dict:
                value = obj.toString(depth + 1)

            case string:
                value = Sprintf("'%v'", obj)

            default:
                value = Sprintf("%+v", obj)
        }

        s += Sprintf("%v  %+v: %+v,\n", space, key, value)
    }

    s += space + "}"

    return s
}

func (self *orderedDict) String() string {
    return self.toString(0)
    // return self.dict.String()
}

func (self *orderedDict) Set(key, value interface{}) {
    if self.keys.Contain(key) == false {
        self.keys.Append(key)
    }

    self.dict[key] = value
}

func (self *orderedDict) Get(key interface{}) interface{} {
    return self.dict[key]
}

func (self *orderedDict) Keys() Array {
    return self.keys
}

func (self *orderedDict) Values() Array {
    values := Array{}

    for _, key := range self.keys {
        values.Append(self.Get(key))
    }

    return values
}

func (self *orderedDict) Items() []orderedDictItem {
    items := []orderedDictItem{}

    for _, key := range self.keys {
        items = append(items, self.Get(key))
    }

    return items
}
