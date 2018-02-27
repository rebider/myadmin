<template>
  <el-menu class="navbar" mode="horizontal">
    <hamburger class="hamburger-container" :toggleClick="toggleSideBar" :isActive="sidebar.opened"></hamburger>
    <breadcrumb></breadcrumb>
    <el-dialog title="修改密码" :visible.sync="dialogFormVisible" width="500px">
      <el-form ref="form" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="旧密码" prop="oldPwd" >
          <el-input v-model="form.oldPwd" type="password"></el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="newPwd">
          <el-input v-model="form.newPwd" type="password"></el-input>
        </el-form-item>
        <el-form-item label="重复新密码" prop="newPwd2">
          <el-input v-model="form.newPwd2" type="password"></el-input>
        </el-form-item>
        <!--<el-form-item>-->
          <!--<el-button type="primary" @click="onSubmit">修改</el-button>-->
        <!--</el-form-item>-->
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="onSubmit">确 定</el-button>
      </div>
    </el-dialog>
    <el-dropdown class="avatar-container" trigger="click">
      <div class="avatar-wrapper">
        <!--<strong></strong> Hi {{name}}-->
        <strong>{{name}}</strong>
        <img class="user-avatar" :src="avatar+'?imageView2/1/w/80/h/80'">
        <!--<img class="user-avatar" :src="avatar+'?imageView2/1/w/80/h/80'">-->
        <!--<i class="el-icon-caret-bottom"></i>-->
      </div>
      <el-dropdown-menu class="user-dropdown" slot="dropdown">
        <router-link class="inlineBlock" to="/">
          <el-dropdown-item>
            主页
          </el-dropdown-item>
        </router-link>
        <el-dropdown-item divided>
          <span @click="dialogFormVisible = true" style="display:block;">修改密码</span>
        </el-dropdown-item>
        <!--<router-link class="inlineBlock" to="/changePasswd">-->
          <!--<el-dropdown-item>-->
            <!--修改密码-->
          <!--</el-dropdown-item>-->
        <!--</router-link>-->
        <el-dropdown-item divided>
          <span @click="logout" style="display:block;">安全退出</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
  </el-menu>


  <!--<el-button type="text" @click="dialogFormVisible = true">打开嵌套表单的 Dialog</el-button>-->
</template>

<script>
import { mapGetters } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'
import { changePasswd } from '@/api/user'
import { Message } from 'element-ui'

export default {
  data() {
    const validatePass = (rule, value, callback) => {
      if (value.length < 5) {
        callback(new Error('密码不能小于5位'))
      } else {
        callback()
      }
    }
    return {
      avatar: 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif',
      dialogFormVisible: false,
      formLabelWidth: '120px',
      form: {
        oldPwd: '',
        newPwd: '',
        newPwd2: ''
      },
      rules: {
        oldPwd: [{ required: true, trigger: 'blur', message: '请输入旧密码' }],
        newPwd: [{ required: true, trigger: 'blur', validator: validatePass }],
        newPwd2: [{ required: true, trigger: 'blur', validator: validatePass }]
      }
    }
  },
  components: {
    Breadcrumb,
    Hamburger
  },
  computed: {
    ...mapGetters([
      'sidebar',
      'name'
      // 'avatar'
    ])
  },
  methods: {
    toggleSideBar() {
      this.$store.dispatch('ToggleSideBar')
    },
    logout() {
      this.$store.dispatch('LogOut').then(() => {
        location.reload() // 为了重新实例化vue-router对象 避免bug
      })
    },
    onSubmit() {
      this.$refs.form.validate(valid => {
        if (valid) {
          if (this.form.newPwd !== this.form.newPwd2) {
            Message({
              message: '两次输入密码不一致!',
              type: 'error',
              duration: 2 * 1000
            })
            return
          }
          changePasswd(this.form).then(respone => {
            Message({
              message: '修改密码成功!',
              type: 'info',
              duration: 2 * 1000
            })
            this.dialogFormVisible = false
            this.form.oldPwd = ''
            this.form.newPwd = ''
            this.form.newPwd2 = ''
          }).catch(error => {
            console.log(error)
          })
        } else {
          return false
        }
      })
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
.navbar {
  height: 50px;
  line-height: 50px;
  border-radius: 0px !important;
  .hamburger-container {
    line-height: 58px;
    height: 50px;
    float: left;
    padding: 0 10px;
  }
  .screenfull {
    position: absolute;
    right: 90px;
    top: 16px;
    color: red;
  }
  .avatar-container {
    height: 50px;
    display: inline-block;
    position: absolute;
    right: 35px;
    .avatar-wrapper {
      cursor: pointer;
      margin-top: 5px;
      position: relative;
      .user-avatar {
        width: 40px;
        height: 40px;
        border-radius: 10px;
      }
      .el-icon-caret-bottom {
        position: absolute;
        right: -20px;
        top: 25px;
        font-size: 12px;
      }
    }
  }
}
</style>

