namespace DevSmtp.Core.Commands
{
    public class ExpnException : Exception
    {
        public ExpnException(string message)
            : base(message)
        {
        }

        public ExpnException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
