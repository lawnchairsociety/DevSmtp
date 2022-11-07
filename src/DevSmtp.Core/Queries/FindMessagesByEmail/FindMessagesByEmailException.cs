namespace DevSmtp.Core.Queries
{
    public class FindMessagesByEmailException : Exception
    {
        public FindMessagesByEmailException(string message)
            : base(message)
        {
        }

        public FindMessagesByEmailException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
