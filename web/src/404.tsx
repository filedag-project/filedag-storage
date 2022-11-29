import { useHistory } from 'react-router';

const NoFoundPage = () => {
  const history = useHistory();
  return (
    <div>
      <button onClick={() => history.push('/')}>Back Home</button>
    </div>
  );
};

export default NoFoundPage;
