package entity

type Meta struct {
	Page      int `json:"page"`
	Paginate  int `json:"paginate"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

func (r *Meta) CountTotalPage(page, paginate, totalData int) {
	r.Page = page
	r.Paginate = paginate

	if totalData == 0 {
		r.TotalPage = 1
		return
	}

	r.TotalData = totalData
	r.TotalPage = totalData / r.Paginate
	if totalData%r.Paginate > 0 {
		r.TotalPage++
	}

	// if totalData == 0, then totalPage should be 1
	if r.TotalPage == 0 {
		r.TotalPage = 1
	}
}
