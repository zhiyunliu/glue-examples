package sqls

const QueryData = `
SELECT id, age, name, address
FROM demotest t where t.id = @{id} &{t.name}  
`

const QueryLikeData = `
SELECT id, age, name, address
FROM demotest  t 
where 1=1  &{t.address like %address%}
`

const QueryLike2Data = `
SELECT id, age, name, address
FROM demotest  t 
where @{t.address like %address%}
`

const QueryCompareData = `
SELECT id, age, name, address
FROM demotest  t 
where 1=1  &{t.age >= age}
`

const QueryCompare2Data = `
SELECT id, age, name, address
FROM demotest  t 
where @{t.age >= age}
`

const QueryInData = `
SELECT id, age, name, address
FROM demotest  t 
where 1=1  &{t.id in ids}
`

const QueryNotInData = `
SELECT id, age, name, address
FROM demotest  t 
where 1=1  &{t.id not in ids}
`

const QueryCTE = `

WITH indexed_ids AS (
	SELECT id
	FROM unnest(@{ids}::int[]) WITH ORDINALITY AS t(id, idx)
	ORDER BY id  
)
SELECT u.*
FROM demotest u
JOIN indexed_ids i ON u.id = i.id
`

const QueryNotIn2Data = `
SELECT id, age, name, address
FROM demotest  t 
where @{t.id not in ids}
`
