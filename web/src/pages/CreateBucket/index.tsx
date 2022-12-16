import classNames from 'classnames';
import styles from './style.module.scss';
import {Button, Form, Input} from 'antd';
import {useNavigate} from 'react-router-dom';
import {RouterPath} from '@/router/RouterConfig';
import bucketsStore from '@/store/modules/buckets';

const CreateBucket = () => {
    const navigate = useNavigate();
    const [form] = Form.useForm();
    
    const create = async () => {
        const bucketName = form.getFieldValue('bucketName1');
        console.log(bucketName,1233);
        try {
            await form.validateFields();
            const path = `/${bucketName}`
            await bucketsStore.fetchCreate(path);
            navigate(RouterPath.buckets);
        } catch (e) {
            
        }
    };

    return (
      <div className={classNames(styles.createBucket)}>
        <Form form={form} autoComplete="off">
            <Form.Item
                name="bucketName1"
                rules={[
                    {required: true, message: 'Please enter bucket name'},
                ]}
                extra="Bucket name must be globally unique and cannot contain spaces or uppercase letters."
            >
                <Input
                    placeholder="Please enter bucket name"
                />
            </Form.Item>

            <Form.Item>
                <Button type="primary" onClick={create}>
                    Create
                </Button>
            </Form.Item>
        </Form>

        <div className={classNames(styles.right)}>
          <div className={classNames(styles.title)}>Buckets</div>
          <div className={classNames(styles.description)}>
            <p>FileDAG uses buckets to organize objects. A bucket is similar to a folder or directory in a filesystem, where each bucket can hold an arbitrary number of objects.</p>
          </div>
        </div>
    </div>
    );
};

export default CreateBucket;
