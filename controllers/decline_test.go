package controllers

import "testing"
type addTest struct{
	a int
	b int 
}
var addTests =[]addTest{
	addTest{2,3},
}

func TestDecline(t *testing.T){

	// for _,test:=range addTests{
	// 	if output :=Decline(test.a);output != test.b{
	// 		t.Errorf("output %q not equal to %q",output,test.b)
	// 	}
	// }
}