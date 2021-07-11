# casbin-query

#### User settings

| User  | Object | Action        |
| :---- | :----- | :------------ |
| david | ewok   | read , update |

#### Role settings

| Role | Object | Action |
| :--- | :----- | :----- |
| A    | ewok   | read   |
| A    | fluffy | read   |
| A    | gordo  | update |
| B    | gordo  | read   |

#### Assigned role

| Role | User    |
| :--- | :------ |
| A    | bob     |
| B    | charlie |

The results when `david` find pets for read is `ewok`.

The results when `charlie` find pets for read is `gordo`.

The results when `bob` find pets for read are `ewok` nad `fluffy`.

#### casbin_rule

| id  | ptype | v0           | v1         | v2     | v3  | v4  | v5  |
| :-- | :---- | :----------- | :--------- | :----- | :-- | :-- | :-- |
| 1   | p     | owner_A      | pet_ewok   | read   |     |     |     |
| 2   | p     | owner_A      | pet_fluffy | read   |     |     |     |
| 3   | p     | owner_A      | pet_gordo  | update |     |     |     |
| 4   | p     | owner_B      | pet_gordo  | read   |     |     |     |
| 5   | p     | user_david   | pet_ewok   | read   |     |     |     |
| 6   | p     | user_david   | pet_fluffy | update |     |     |     |
| 7   | g     | user_bob     | owner_A    |        |     |     |     |
| 8   | g     | user_charlie | owner_B    |        |     |     |     |

#### pet

| id  | version | created_at          | updated_at          | name   |
| :-- | :------ | :------------------ | :------------------ | :----- |
| 1   | 1       | 2021-07-11 14:13:43 | 2021-07-11 14:13:43 | ewok   |
| 2   | 1       | 2021-07-11 14:13:43 | 2021-07-11 14:13:43 | fluffy |
| 3   | 1       | 2021-07-11 14:13:43 | 2021-07-11 14:13:43 | gordo  |

#### SQL

```
SELECT SUBSTRING_INDEX(tp.v1, '_', -1) AS name
FROM casbin_rule tg
INNER JOIN casbin_rule tp ON tg.v1 = tp.v0
WHERE tg.v0 = 'user_david'
AND tg.ptype = 'g'
AND tp.ptype = 'p'
AND tp.v2 = 'read'

UNION

SELECT SUBSTRING_INDEX(tp.v1, '_', -1) AS name
FROM casbin_rule tp
WHERE tp.v0 = 'user_david'
AND tp.ptype = 'p'
AND tp.v2 = 'read'
) AS t3 ON `pet`.`name`= t3.name ORDER BY `pet`.`name`
```
