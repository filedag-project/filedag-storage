import styles from './action.module.scss';
import { Button } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { RouterPath } from '@/router/RouterConfig';
import { useHistory } from 'react-router';


const Action = () => {
  const history = useHistory();
  //const onSearch = () => {};
  const create = ()=>{
    history.push(RouterPath.createBucket);
  }
  return (
    <div className={styles.action}>
      {/* <div className={styles.search}>
        <Search
          placeholder="input search text"
          allowClear
          onSearch={onSearch}
          enterButton="Search"
        />
      </div> */}
      <div className={styles.operation}>
        <Button type="primary" icon={<PlusOutlined />} onClick={create}>
          Create
        </Button>
      </div>
    </div>
  );
};

export default Action;
