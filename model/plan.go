package model

type Plan struct {
	Project 	string		`json:"project"`
	Links   	[]string	`json:"links"`

	LinkList    []Link
	CurrentQps  map[string]int
}

 func (p *Plan) Merge(project *Project) {
	for _, s := range p.Links {
		for _, link := range project.Links {
			if s == link.Name {
				for _, s := range project.Services{
					if s.Name == link.ServiceName {
						link.Service = s
					}
				}
				p.LinkList = append(p.LinkList, link)
				p.CurrentQps[link.Name] = link.TargetQps  // ?
			}
		}
	}
}