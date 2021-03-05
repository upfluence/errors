package tags

import "github.com/upfluence/errors/base"

func GetTags(err error) map[string]interface{} {
	var tags map[string]interface{}

	for {
		if err == nil {
			break
		}

		if t, ok := err.(interface{ Tags() map[string]interface{} }); ok {
			ts := t.Tags()

			if len(ts) > 0 && tags == nil {
				tags = make(map[string]interface{}, len(ts))
			}

			for k, v := range ts {
				if _, ok := tags[k]; !ok {
					tags[k] = v
				}
			}
		}

		err = base.UnwrapOnce(err)
	}

	return tags
}
