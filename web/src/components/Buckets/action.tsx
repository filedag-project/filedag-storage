import styles from './action.module.scss';
import { Button } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { RouterPath } from '@/router/RouterConfig';
import { useNavigate } from 'react-router-dom';


const Action = () => {
  const navigate = useNavigate();
  const create = ()=>{
    navigate(RouterPath.createBucket);
  }
  return (
    <div className={styles.action}>
      <div className={styles.operation}>
        <Button className='bg-btn' type="primary" icon={<PlusOutlined />} onClick={create}>
          Create
        </Button>
      </div>
    </div>
  );
};

export default Action;
