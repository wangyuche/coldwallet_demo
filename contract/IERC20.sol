// SPDX-License-Identifier: MIT

pragma solidity ^0.8.4;

/**
 * @dev ERC20 接口合約.
 */
interface IERC20 {
    /**
     * @dev 釋放條件：當 `value` 單位的貨幣從帳戶 (`from`) 轉帳到另一帳戶 (`to`)時.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev 釋放條件：當 `value` 單位的貨幣從帳戶 (`owner`) 授權給另一帳戶 (`spender`)時.
     */
    event Approval(address indexed owner, address indexed spender, uint256 value);

    /**
     * @dev 返回代幣總供給.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev 返回帳戶`account`所持有的代幣數.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev 轉帳 `amount` 單位代幣，從調用者帳戶到另一帳戶 `to`.
     *
     * 如果成功，返回 `true`.
     *
     * 釋放 {Transfer} 事件.
     */
    function transfer(address to, uint256 amount) external returns (bool);

    /**
     * @dev 返回`owner`帳戶授權給`spender`帳戶的額度，默認為0。
     *
     * 當{approve} 或 {transferFrom} 被調用時，`allowance`會改變.
     */
    function allowance(address owner, address spender) external view returns (uint256);

    /**
     * @dev 調用者帳戶給`spender`帳戶授權 `amount`數量代幣。
     *
     * 如果成功，返回 `true`.
     *
     * 釋放 {Approval} 事件.
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev 通過授權機制，從`from`帳戶向`to`帳戶轉帳`amount`數量代幣。轉帳的部分會從調用者的`allowance`中扣除。
     *
     * 如果成功，返回 `true`.
     *
     * 釋放 {Transfer} 事件.
     */
    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) external returns (bool);
}