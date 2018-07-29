package test

import(
    "github.com/agiledragon/trans-dsl"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "github.com/agiledragon/trans-dsl/test/context/action"
    "github.com/agiledragon/trans-dsl/test/context"
    "github.com/agiledragon/trans-dsl/test/context/spec"
)

func newIfTrans() *transdsl.Transaction {
    trans := &transdsl.Transaction {
        Fragments: []transdsl.Fragment{
            &transdsl.Optional{
                Spec: new(spec.IsAbcExist),
                IfFrag: new(action.StubConnectAbc),
            },
            new(action.StubActivateSomething),
        },
    }
    return trans
}

func newElseTrans() *transdsl.Transaction {
    trans := &transdsl.Transaction {
        Fragments: []transdsl.Fragment{
            &transdsl.Optional{
                Spec: new(spec.IsAbcExist),
                IfFrag: new(action.StubConnectAbc),
                ElseFrag: new(action.StubConnectDef),
            },
            new(action.StubActivateSomething),
        },
    }
    return trans
}

func TestIfTrans(t *testing.T) {
    trans := newIfTrans()
    Convey("TestIfTrans", t, func() {

        Convey("trans exec succ when spec is true", func() {
            transInfo := &transdsl.TransInfo{
                AppInfo: &context.StubInfo{
                    X: "abc",
                    Y: 1,
                },
            }
            err := trans.Exec(transInfo)
            So(err, ShouldEqual, nil)
            So(transInfo.AppInfo.(*context.StubInfo).Y, ShouldEqual, 2)
        })

        Convey("trans exec succ when spec is false", func() {
            transInfo := &transdsl.TransInfo{
                AppInfo: &context.StubInfo{
                    X: "def",
                    Y: 1,
                },
            }
            err := trans.Exec(transInfo)
            So(err, ShouldEqual, nil)
            So(transInfo.AppInfo.(*context.StubInfo).Y, ShouldEqual, 1)
        })

        Convey("ifspec rollback", func() {
            transInfo := &transdsl.TransInfo{
                AppInfo: &context.StubInfo{
                    X: "test",
                    Y: 1,
                },
            }
            err := trans.Exec(transInfo)
            So(err, ShouldNotEqual, nil)
            So(transInfo.AppInfo.(*context.StubInfo).Y, ShouldEqual, 0)
        })
    })
}

func TestElseTrans(t *testing.T) {
    trans := newElseTrans()
    Convey("TestElseTrans", t, func() {

        Convey("trans exec succ when spec is false", func() {
            transInfo := &transdsl.TransInfo{
                AppInfo: &context.StubInfo{
                    X: "def",
                    Y: 1,
                },
            }
            err := trans.Exec(transInfo)
            So(err, ShouldEqual, nil)
            So(transInfo.AppInfo.(*context.StubInfo).Y, ShouldEqual, 3)
        })


        Convey("elsespec rollback", func() {
            transInfo := &transdsl.TransInfo{
                AppInfo: &context.StubInfo{
                    X: "test",
                    Y: 1,
                },
            }
            err := trans.Exec(transInfo)
            So(err, ShouldNotEqual, nil)
            So(transInfo.AppInfo.(*context.StubInfo).Y, ShouldEqual, -1)
        })
    })
}

