<template>

    <imp-panel>
        <h3 class="box-title" slot="header" style="width: 100%;">
            <el-row style="width: 100%;">
                <el-col :span="12">
                    <el-button class="filter-item" style="margin-left: 10px;" @click="handleCreate" type="primary"
                               icon="el-icon-edit">新增
                    </el-button>
                    <!--<router-link :to="{ path: 'userAdd'}">-->
                    <!--<el-button type="primary" icon="plus">新增</el-button>-->
                    <!--</router-link>-->
                </el-col>
            </el-row>
        </h3>
        <div slot="body">
            <el-table
                    :data="tableData.rows"
                    border
                    highlight-current-row
                    stripe
                    style="width: 100%"
                    v-loading="listLoading"
                    @selection-change="handleSelectionChange">
                <!--checkbox 适当加宽，否则IE下面有省略号 https://github.com/ElemeFE/element/issues/1563-->
                <el-table-column
                        type="selection"
                        width="50">
                </el-table-column>
                <el-table-column
                        prop="Id"
                        label="ID"
                        width="100">
                </el-table-column>
                <!--<el-table-column-->
                <!--label="照片" width="76">-->
                <!--<template slot-scope="scope">-->
                <!--<img :src='scope.row.avatar' style="height: 35px;vertical-align: middle;" alt="">-->
                <!--</template>-->
                <!--</el-table-column>-->
                <el-table-column
                        prop="Account"
                        width="200"
                        label="帐号">
                </el-table-column>
                <el-table-column
                        prop="Name"
                        width="200"
                        label="名称">
                </el-table-column>
                <el-table-column
                        prop="LoginTimes"
                        width="100"
                        label="登录次数">
                </el-table-column>
                <el-table-column
                        prop="LastLoginTime"
                        width="150"
                        label="最后登录时间">
                </el-table-column>
                <el-table-column
                        prop="LastLoginIp"
                        width="150"
                        label="最后登录IP">
                </el-table-column>
                <!--<el-table-column-->
                <!--prop="email"-->
                <!--label="邮箱">-->
                <!--</el-table-column>-->
                <el-table-column width="150" label="状态">

                    <template slot-scope="scope">
                        <el-tag :type="scope.row.Status | statusFilter">{{scope.row.Status | statusNameFilter }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column
                        prop="RoleIds"
                        width="300"
                        label="角色">
                </el-table-column>
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-button
                                size="small"
                                type="primary"
                                icon="edit"
                                @click="handleEdit(scope.row)">编辑
                        </el-button>
                        <el-button
                                size="small"
                                type="danger"
                                @click="">删除
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>

            <el-pagination
                    background
                    @size-change="handleSizeChange"
                    @current-change="handleCurrentChange"
                    :current-page="tableData.pagination.pageNo"
                    :page-sizes="[5, 10, 20]"
                    :page-size="tableData.pagination.pageSize"
                    layout="total, sizes, prev, pager, next"
                    :total="tableData.pagination.total">
            </el-pagination>

            <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible" width="500px">
                <el-form ref="form" :model="form" :rules="rules" label-width="120px">
                    <el-form-item label="帐号" prop="account">
                        <el-input v-model="form.account" type="text" :disabled="isAccountDisabled"></el-input>
                    </el-form-item>
                    <el-form-item label="名称" prop="name">
                        <el-input v-model="form.name" type="text"></el-input>
                    </el-form-item>
                    <el-form-item label="状态" prop="status">
                        <el-radio-group v-model="form.status">
                            <el-radio :label="0">禁用</el-radio>
                            <el-radio :label="1">启用</el-radio>
                        </el-radio-group>
                        <!--<el-input v-model="form.status" type="text"></el-input>-->
                    </el-form-item>
                    <el-form-item label="角色" prop="roleIds">
                        <el-select v-model="form.roleIds" multiple placeholder="请选择">
                            <el-option
                                    v-for="item in roleList"
                                    :key="item.Id"
                                    :label="item.Name"
                                    :value="item.Id">
                            </el-option>
                        </el-select>

                        <!--<el-input v-model="form.roles" type="text"></el-input>-->
                    </el-form-item>
                    <el-form-item label="密码" prop="passwd">
                        <el-input v-model="form.passwd" type="text"></el-input>
                    </el-form-item>
                </el-form>
                <div slot="footer" class="dialog-footer">
                    <el-button @click="dialogFormVisible = false">取 消</el-button>
                    <el-button type="primary" @click="onSubmit">保 存</el-button>
                </div>
            </el-dialog>

        </div>


    </imp-panel>
</template>

<script>
  import panel from '@/components/panel.vue'
  import { userList, edit } from '@/api/user'
  import { roleList } from '@/api/role'
  import { Message } from 'element-ui'
  // import * as api from "../../api"
  // import testData from "../../../static/data/data.json"
  // import * as sysApi from '../../services/sys'

  export default {
    components: {
      'imp-panel': panel
    },
    data() {
      return {
        currentRow: {},
        dialogVisible: false,
        dialogLoading: false,
        defaultProps: {
          children: 'children',
          label: 'name',
          id: 'id'
        },
        roleList: [],
        roleTree: [],
        listLoading: false,
        searchKey: '',
        tableData: {
          pagination: {
            total: 0,
            pageNo: 1,
            pageSize: 10,
            parentId: 0
          }
        },
        rows: [],
        dialogFormVisible: false,
        formLabelWidth: '120px',
        form: {
          id: '',
          account: '',
          name: '',
          status: '1',
          roleIds: '',
          passwd: ''
        },
        isAccountDisabled: false,
        rules: {
          account: [{ required: true, message: '请输入帐号', trigger: 'blur' }],
          name: [{ required: true, trigger: 'blur', message: '请输入昵称' }]
          // status: [{ required: true, trigger: 'blur' }]
        },
        dialogStatus: '',
        textMap: {
          edit: '编辑用户',
          create: '新增用户'
        },
        options: [
          { value: 'admin', label: '超级管理员' },
          { value: 'yanfa', label: '研发' }
        ]
      }
    },
    filters: {
      statusFilter(status) {
        const statusMap = {
          '0': 'danger',
          '1': 'success'
        }
        return statusMap[status]
      },
      statusNameFilter(status) {
        const statusNameMap = {
          '0': '禁用',
          '1': '启用'
        }
        return statusNameMap[status]
      }
    },
    methods: {
      handleSelectionChange(val) {
      },
      handleSizeChange(val) {
        this.tableData.pagination.pageSize = val
        this.loadData()
      },
      handleCurrentChange(val) {
        this.tableData.pagination.pageNo = val
        this.loadData()
      },
      loadData() {
        userList({
          key: this.searchKey,
          pageSize: this.tableData.pagination.pageSize,
          pageNo: this.tableData.pagination.pageNo
        }).then(res => {
          const data = res.data
          console.log(data)
          console.log(data.total)
          console.log(data.rows)
          this.tableData.rows = data.rows
          this.tableData.pagination.total = data.total
        }).catch(error => {
          console.log(error)
        })
        roleList({
        }).then(res => {
          const data = res.data
          console.log(data)
          console.log(data.total)
          console.log(data.rows)
          this.roleList = data.rows
          // this.tableData.pagination.total = data.total
        }).catch(error => {
          console.log(error)
        })
      },
      handleCreate() {
        console.log(666666666666666)
        this.dialogStatus = 'create'
        this.dialogFormVisible = true
        this.isAccountDisabled = false
        this.form.id = ''
        this.form.account = ''
        this.form.name = ''
        this.form.status = ''
        this.form.roleIds = ''
        // this.$nextTick(() => {
        //     this.$refs['dataForm'].clearValidate()
        // })
      },
      handleEdit(row) {
        console.log(7777777777)
        console.log(row)
        console.log(row.Account)
        console.log(row.RoleIds)
        this.dialogStatus = 'edit'
        this.dialogFormVisible = true
        this.isAccountDisabled = true
        this.form.id = row.Id
        this.form.account = row.Account
        this.form.name = row.Name
        this.form.status = row.Status
        this.form.roleIds = row.RoleIds
        // this.$nextTick(() => {
        //     this.$refs['dataForm'].clearValidate()
        // })
      },
      onSubmit() {
        this.$refs.form.validate(valid => {
          if (valid) {
            edit(this.form).then(respone => {
              Message({
                message: '编辑用户成功!',
                type: 'info',
                duration: 2 * 1000
              })
              this.dialogFormVisible = false
              this.loadData()
            }).catch(error => {
              Message({
                message: '编辑用户失败!' + error,
                type: 'info',
                duration: 2 * 1000
              })
            })
          } else {
            return false
          }
        })
      }
    },
    created() {
      this.loadData()
    }
  }
</script>
<style>
    .el-pagination {
        float: right;
        margin-top: 15px;
    }
</style>
