import { addUserType, changePasswordType, statusType, userStatusType } from "@/models/DashboardModel";
import { tokenType } from "@/models/RouteModel";
import userStore from "@/store/modules/user";
import { Cookies, SESSION_TOKEN } from "@/utils/cookies";
import { DeleteOutlined, PlusOutlined, SafetyOutlined } from "@ant-design/icons";
import { Button, Form, Input,Switch, Modal, Table, Tooltip, notification } from "antd";
import { observer } from "mobx-react";
import { useEffect, useState } from "react";
import jwt_decode from "jwt-decode";
import styles from './style.module.scss';

const User = (props:any) => {
  const [addUserForm] = Form.useForm();
  const [changePassWordForm] = Form.useForm();
  const [deleteUserShow,SetDeleteUserShow] = useState(false);
  const [changePasswordShow,SetChangePasswordShow] = useState(false);
  const [userStatus,SetUserStatusShow] = useState(false);
  const [accessKey,SetAccessKey] = useState('');
  const [defaultStatus,SetDefaultStatus] = useState(false);
  const [status,SetStatus] = useState(false);
  const [addUserShow,SetAddUserShow] = useState(false);
  const [admin,setAdmin] = useState(false);
  const [adminName,setAdminName] = useState('');
  useEffect(()=>{
    const _jwt = Cookies.getKey(SESSION_TOKEN);
    if(_jwt){
      const _token:tokenType = jwt_decode(_jwt);
      const { isAdmin = false,parent }=_token;
      setAdmin(isAdmin);
      setAdminName(parent);
    }
    userStore.fetchUserInfos();
    
  },[]);

  const columns = [
    {
      title: 'Account Name',
      dataIndex: 'account_name',
      key: 'account_name',
    },
    {
      title: 'total Storage Capacity',
      dataIndex: 'total_storage_capacity',
      key: 'total_storage_capacity',
    },
    {
      title: 'Use Storage Capacity',
      dataIndex: 'use_storage_capacity',
      key: 'use_storage_capacity',
    },
    {
      title: 'Bucket',
      dataIndex: 'bucket_infos',
      key: 'bucket_infos',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render:(r,record)=>{
        const _value = r === statusType.on;
        return admin && (record.account_name === adminName)?<></>:<div className='switch-wrap'>
          <Switch checkedChildren={statusType.on} unCheckedChildren={statusType.off} defaultChecked={_value}></Switch>
          <div className='mask'  onClick={()=>{
            SetAccessKey(record.account_name);
            SetUserStatusShow(true);
            SetDefaultStatus(_value);
          }}></div>
        </div>
      }
    },
    {
      title: 'Action',
      key: 'action',
      render: (_, record) => (
        <div className='row-action'>
          {
            admin && (record.account_name === adminName)?<></>:<span onClick={()=>{
              SetAccessKey(record.account_name);
              SetChangePasswordShow(true);
            }}>
              <Tooltip title="Change Password"><SafetyOutlined /></Tooltip>
            </span>
          }
          {
            admin && (record.account_name === adminName)?<></>:<span onClick={()=>{
              SetAccessKey(record.account_name);
              SetDeleteUserShow(true);
            }}>
              <Tooltip title="Delete User"><DeleteOutlined /></Tooltip>
            </span>
          }
          
        </div>
      ),
    },
  ];
  const addUser = async ()=>{
    try {
      await addUserForm.validateFields();
      const username = addUserForm.getFieldValue('username');
      const password = addUserForm.getFieldValue('password');
      const capacity = addUserForm.getFieldValue('capacity');
      const user:addUserType = {
        accessKey:username,
        secretKey:password,
        capacity
      }
      userStore.fetchAddUser(user).then(res=>{
        SetAddUserShow(false);
        userStore.fetchUserInfos();
      })
    } catch (e) {
        
    }
  };

  const deleteUser = ()=>{
    userStore.fetchDeleteUser(accessKey).then(res=>{
      userStore.fetchUserInfos();
      SetDeleteUserShow(false);
    });
  };
  const changePassword = async ()=>{
    await changePassWordForm.validateFields();
    const newSecretKey = changePassWordForm.getFieldValue('newSecretKey');
    const params:changePasswordType = {    
      newSecretKey,
      accessKey
    }
    userStore.fetchChangeUserPassword(params).then(res=>{
      SetChangePasswordShow(false);
      changePassWordForm.resetFields();
      notification.open({
        message: 'Change success',
        description: '',
      });
    });
  };

  const setUserStatus = ()=>{
    const params:userStatusType ={
      accessKey,
      status: status ? statusType.on : statusType.off
    }
    userStore.fetchSetUserStatus(params).then(async res=>{
      SetUserStatusShow(false);
      userStore.fetchUserInfos();
    })
  };

  return <div className={styles.user}>
    <div className={styles.userList}>
      <div className={styles.action}>
        <Button className="bg-btn" type="primary" icon={<PlusOutlined />} onClick={()=>{ SetAddUserShow(true) }}>
          Add User
        </Button>
      </div>
      <Table columns={columns} dataSource={userStore.userInfos} rowKey={(record) => record.account_name + record.status } pagination={false}/>
    </div>
    <Modal
      title="Add User"
      open={addUserShow}
      onOk={addUser}
      onCancel={()=>{ 
        SetAddUserShow(false);
        addUserForm.resetFields();
      }}
      okText="Confirm"
      cancelText="Cancel"
    >
      <Form form={addUserForm} autoComplete="off">
            <Form.Item
                name="username"
                rules={[
                    {required: true, message: 'Please input username'},
                ]}
            >
                <Input
                    placeholder="please enter username"
                />
            </Form.Item>
            <Form.Item
                name="password"
                rules={[
                    {required: true, message: 'Please input password'},
                ]}
            >
                <Input
                    placeholder="please enter password"
                />
            </Form.Item>
            <Form.Item
                name="capacity"
                rules={[
                    {required: true, message: 'Please input capacity'},
                ]}
            >
                <Input
                    placeholder="please enter capacity"
                />
            </Form.Item>
        </Form>
    </Modal>

    <Modal
      title="Delete"
      open={deleteUserShow}
      onOk={deleteUser}
      onCancel={()=>{ SetDeleteUserShow(false) }}
      okText="Confirm"
      cancelText="Cancel"
    >
      <p>Are you sure to delete this data？</p>
    </Modal>

    <Modal
      title="Set Status"
      open={userStatus}
      onOk={setUserStatus}
      onCancel={()=>{ SetUserStatusShow(false) }}
      okText="Confirm"
      cancelText="Cancel"
    >
      <div style={{display:'flex',}}>
        <div className="label" style={{margin:'0 8px 0 0',}}>Status</div>：
        <Switch checkedChildren={statusType.on} unCheckedChildren={statusType.off} defaultChecked={defaultStatus} onChange={(checked)=>{
          SetStatus(checked);
        }}></Switch>
      </div>
    </Modal>
    
    <Modal
      title="Change Password"
      open={changePasswordShow}
      onOk={changePassword}
      onCancel={()=>{
        SetChangePasswordShow(false);
        changePassWordForm.resetFields();
      }}
      okText="Confirm"
      cancelText="Cancel"
    >
      <Form form={changePassWordForm} autoComplete="off">
            <Form.Item
                name="newSecretKey"
                rules={[
                    {required: true, message: 'Please input new password'},
                    { min:8, message: 'password must be at least 8 characters'},
                ]}
            >
                <Input.Password
                    placeholder="please enter new password"
                />
            </Form.Item>
        </Form>
    </Modal>
  </div>;
};

export default observer(User);

