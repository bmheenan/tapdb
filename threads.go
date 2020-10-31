package tapdb

/*
	WITH        RECURSIVE descendants (child, parent) AS
	            (
	            SELECT child
	              ,    parent
	            FROM   threads_parent_child
	            WHERE  parent = '%v'
	            UNION ALL
	            SELECT t.child
	              ,    t.parent
	            FROM   threads_parent_child t
	            JOIN   descendants
	              ON   t.parent = descendants.child
		        )
	SELECT      d.child
	  ,         t.owner
	  ,         (s.stakeholder = '%v') AS tracked
	  , 	    t.costdirect
	  , 	    t.iteration
	FROM        descendants d
	  LEFT JOIN (
			    SELECT thread
			      ,    stakeholder
				FROM   threads_stakeholders
				WHERE  stakeholder = '%v'
	            ) s
	  ON        s.thread = d.child
	  JOIN      threads t
	  ON        t.id = d.child
	ORDER BY    s.stakeholder
*/
