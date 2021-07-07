select c.column_name as colname,
    c.udt_name as coltype,
    c.character_maximum_length as collength,
    c.numeric_scale as precision,
    case c.is_nullable
        when 'YES' then cast('t' as bool)
        when 'NO' then cast('f' as bool)
    end as nullable
from information_schema.columns c
where c.table_name = $1;