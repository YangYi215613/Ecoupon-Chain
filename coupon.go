package main

// 卡券数据
type Coupon struct {
	Data []byte  // 原始数据
}

func NewCoupon() *Coupon {
	return &Coupon{
		Data: []byte("电子卡券数据"),
	}
}