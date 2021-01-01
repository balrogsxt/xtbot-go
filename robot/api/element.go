package api

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
)

type IMsg interface {
}



//目前封装支持的消息结构
type (
	//文字消息
	Text struct {
		Content string //文字内容
	}

	//图片消息
	Image struct {
		Id string //所处在腾讯的图片ID
	}
	Face struct {
		Name string //显示名称
		Id int32 //faceId
	}

	At struct {
		QQ int64 //At的QQ号码,为0时为@All
		Display string //显示的名称
	}
	RedPack struct {
		Title string //红包名称
	}
)
//解析消息为字符串
func ToString(list []IMsg) string {
	str := ""
	for _, elem := range list {
		switch e := elem.(type) {
		case Text:
			str += e.Content
			break
		case Image:
			str += fmt.Sprintf("[Image:%s]", e.Id)
			break
		case Face:
			str += fmt.Sprintf("[%s]", e.Name)
			break
		case At:
			str += e.Display
			break
		case RedPack:
			str += fmt.Sprintf("[QQ红包:%s]", e.Title)
			break
		}
	}
	return str
}

//解析封装消息->原始消息
func ParseToOldElement(elements []IMsg) []message.IMessageElement{
	list := make([]message.IMessageElement,0)
	for _, elem := range elements {
		switch e := elem.(type) {
		case Text:
			list = append(list,message.NewText(e.Content))
			break
		case Image:
			list = append(list,&message.ImageElement{
				Filename: e.Id,
			})
			break
		case Face:
			list = append(list,message.NewFace(e.Id))
			break
		case At:
			list = append(list,message.NewAt(e.QQ,e.Display))
			break
		}
	}
	return list
}
//解析原始消息->封装消息
func ParseElement(elements []message.IMessageElement) []IMsg {
	list := make([]IMsg,0)

	for _, elem := range elements {
		switch e := elem.(type) {
		case *message.TextElement:
			list = append(list, Text{
				Content: e.Content,
			})
			break
		case *message.ImageElement:
			list = append(list, Image{
				Id: e.Filename,
			})
		break
			case *message.FaceElement:
				list = append(list, Face{
					Name:e.Name,
					Id:e.Index,
				})
			break
		case *message.GroupImageElement:
			list = append(list, Image{
				Id: e.ImageId,
			})
			break
		case *message.AtElement:
			list = append(list, At{
				QQ: e.Target,
				Display: e.Display,
			})
			break
		case *message.RedBagElement:
			list = append(list, RedPack{
				Title: e.Title,
			})
			break
		//case *ReplyElement:
		//	res += "[Reply:" + strconv.FormatInt(int64(e.ReplySeq), 10) + "]"
		}
	}
	return list
}


//消息Id结构
type MsgId struct {
	Id int32 //消息Id
	InternalId int32 //内部消息ID
}
//群组消息Id
type GroupMsgId struct {
	MsgId
	Group
}
//撤回消息
func (this *GroupMsgId) RecallMessage() {
	RecallGroupMessage(this.Group.Id,this.MsgId)
}


//QQ用户结构
type QQ struct {
	Uin int64 //QQ号码
	Name string //QQ昵称
}



//群组用户
type GroupUser struct {
	QQ
	CardName string //群名片
}

//获取显示名称
func (this *GroupUser) GetDisplayName() string {
	if len(this.CardName) == 0 {
		return this.Name
	}else{
		return this.CardName
	}
}
//群组结构
type Group struct {
	Id int64 //群组Id
	Name string //群组名称
}
//发送群组消息
func (this *Group) SendGroupMessage(list []IMsg) GroupMsgId {
	return SendGroupMessage(this.Id,list)
}
//撤回群组消息
func (this *Group) RecallGroupMessage(msgId MsgId) {
	RecallGroupMessage(this.Id,msgId)
}
//快速获取图片数据
func (this *Group) NewImage(id_path_url string) IMsg {
	return NewImage(this.Id,id_path_url)
}

