import React from 'react';
import { Segment, Header } from 'semantic-ui-react';

const About = () => (
  <>
    <Segment>
      <Header as="h3">说明</Header>
      <hr></hr>
      <h4>API</h4>
      <h5>获取 Access Token</h5>
      <p>请求方法：GET
        URL：/api/wechat/access_token
        无参数，但是需要设置 HTTP 头部：Authorization: 【token】
      </p>
      <h5>通过验证码查询用户 ID</h5>
      <p>请求方法：GET
        URL：/api/wechat/user?code=【code】
        需要设置 HTTP 头部：Authorization: 【token】
        注意
        需要将 【token】和 【code】替换为实际的内容。
      </p>
      <hr></hr>
      <br></br>
      <p> 加* 的菜单为管理员使用</p>
      <br></br>
      <hr></hr>
      GitHub:{' '}
      <a href="https://github.com/songquanpeng/react-template">
        https://github.com/songquanpeng/react-template
      </a>
    </Segment>
  </>
);

export default About;
