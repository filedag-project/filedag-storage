import classNames from 'classnames';
import styles from './style.module.scss';
import {useNavigate} from 'react-router-dom';
import {RouterPath} from '@/router/RouterConfig';
import { Button } from 'antd';

const NotFound = () => {
    const navigate = useNavigate();
    const back = ()=>{
        navigate(RouterPath.home)
    }
    return (
        <div className={classNames(styles.notFound)}>
            <div className={styles.tips}>Page Not Found</div>
            <Button className='bg-btn' type="primary" onClick={back}>Back</Button>
        </div>
    );
};

export default NotFound;
