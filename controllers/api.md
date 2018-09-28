## API list
### 获取老师列表
* GET
* http://localhost:8090/v1/dean/vote/{{voteCode}}, {{voteCode}}就是投票码参数
* 返回值
	{
	    "code": 0,
	    "msg": "success",
	    "data": {
	        "Grade": 1,
	        "Index": 3,
	        "ID": 259,
	        "Name": "高一三班",
	        "Teachers": [
	            {
	                "ID": 1,
	                "Gender": 1,
	                "Name": "辛俊任",
	                "Mobile": "12345678910"
	            },
	            {
	                "ID": 2,
	                "Gender": 2,
	                "Name": "戴茜",
	                "Mobile": "12346543728"
	            }
	        ]
	    }
	}


### 投票接口
* POST
* http://localhost:8090/v1/dean/vote
	+ body 参数
	{
		"VoteCode":"aronic",
		"Scores":[{
			"TeacherID":1,
			"Score":1
		},{
			"TeacherID":2,
			"Score":1
		},{
			"TeacherID":3,
			"Score":1
		}]
	}
* 返回
	{
	    "code": 0,
	    "msg": "success",
	    "data": null
	}
