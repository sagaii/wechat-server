import React, { useEffect, useState } from 'react';
import { Form, Grid } from 'semantic-ui-react';
import { API, showError } from '../helpers';

const LogViewTable = (props) => {
  const [logText, setLogText] = useState(''); // 日志内容的状态变量
  const [loading, setLoading] = useState(false); // 加载状态

  const fetchLog = async (logType) => {
    setLoading(true);
    try {
      const res = await API.get(`/api/logs/${logType}`);
      const { success, message, data } = res.data;
      console.log("logType",logType,data);
      if (success) {
        // data 是一个字符串，直接显示在 textarea 中
        setLogText(data);
      } else {
        showError(message);
      }
    } catch (error) {
      showError('网络或服务器错误');
    }
    setLoading(false);
  };

  useEffect(() => {
    if (props.type === 'error') {
      fetchLog('error');
    } else if (props.type === 'common') {
      fetchLog('common');
    }
  }, [props.type]); // 依赖 props.type 以便于切换日志类型时重新加载

  return (
    <Grid columns={1}>
      <Grid.Column>
        <Form loading={loading}>
          <Form.Group widths="equal">
            <Form.TextArea
              label={
                <p>                 
                </p>
              }
              placeholder="日志内容"
              value={logText}
              name="logTextInfo"
              style={{ minHeight: 300, fontFamily: 'JetBrains Mono, Consolas' }}
              readOnly
            />
          </Form.Group>
        </Form>
      </Grid.Column>
    </Grid>
  );
};

export default LogViewTable;
